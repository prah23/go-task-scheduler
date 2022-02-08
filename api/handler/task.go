package handler

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"path/filepath"
	"strconv"
	"sync"
	"task-scheduler/api/views"
	"task-scheduler/pkg"
	"task-scheduler/pkg/task"
)

func create(sv task.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			views.Error(pkg.ErrMethodNotAllowed.Error(), http.StatusMethodNotAllowed, w)
			return
		}

		var newTaskRequest task.TaskCreateRequest
		err := json.NewDecoder(r.Body).Decode(&newTaskRequest)
		if err != nil {
			views.Error(err.Error(), http.StatusBadRequest, w)
			return
		}

		newTask := &task.Task{
			ID:          strconv.Itoa(rand.Int()),
			Name:        newTaskRequest.Name,
			Description: newTaskRequest.Description,
			Timeout:     newTaskRequest.Timeout,
		}

		t, err := sv.CreateTask(r.Context(), newTask)
		if err != nil {
			views.Error(err.Error(), http.StatusInternalServerError, w)
			return
		}

		views.Success(t, http.StatusCreated, w)
	})
}

func submit(sv task.Service, wg *sync.WaitGroup) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			views.Error(pkg.ErrMethodNotAllowed.Error(), http.StatusMethodNotAllowed, w)
			return
		}

		taskID := r.FormValue("id")

		isTaskRunning, err := sv.IsTaskRunning(r.Context(), taskID)
		if err != nil {
			views.Error("task does not exist", http.StatusBadRequest, w)
			return
		}

		if isTaskRunning {
			views.Error("task is already running", http.StatusBadRequest, w)
			return
		}

		taskFile, taskFileHeader, err := r.FormFile("script")
		if err != nil {
			views.Error(err.Error(), http.StatusInternalServerError, w)
			return
		}

		ext := filepath.Ext(taskFileHeader.Filename)
		if ext != ".sh" {
			views.Error("only shell scripts are accepted", http.StatusBadRequest, w)
			return
		}

		contents, err := ioutil.ReadAll(taskFile)
		if err != nil {
			views.Error(err.Error(), http.StatusInternalServerError, w)
			return
		}

		_, err = task.StoreScriptInFile(contents, taskID)
		if err != nil {
			views.Error(err.Error(), http.StatusInternalServerError, w)
			return
		}

		_, err = sv.RunTask("sh", r.Context(), wg, taskID)
		if err != nil {
			views.Error(err.Error(), http.StatusInternalServerError, w)
			return
		}

		views.Success("task started successfully", http.StatusAccepted, w)
	})
}

func viewAll(sv task.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			views.Error(pkg.ErrMethodNotAllowed.Error(), http.StatusMethodNotAllowed, w)
			return
		}

		tasks, err := sv.GetAllTasks(r.Context())
		if err != nil {
			views.Error(err.Error(), http.StatusInternalServerError, w)
			return
		}
		views.Success(tasks, http.StatusOK, w)
	})
}

func view(sv task.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			views.Error(pkg.ErrMethodNotAllowed.Error(), http.StatusMethodNotAllowed, w)
			return
		}

		var newTask task.TaskIDRequest
		err := json.NewDecoder(r.Body).Decode(&newTask)
		if err != nil {
			views.Error(err.Error(), http.StatusBadRequest, w)
			return
		}

		task, err := sv.GetTask(r.Context(), newTask.ID)
		if err != nil {
			views.Error(err.Error(), http.StatusBadRequest, w)
			return
		}
		views.Success(task, http.StatusOK, w)
	})
}

func cancel(sv task.Service, wg *sync.WaitGroup) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			views.Error(pkg.ErrMethodNotAllowed.Error(), http.StatusMethodNotAllowed, w)
			return
		}

		var existingTask task.TaskIDRequest
		err := json.NewDecoder(r.Body).Decode(&existingTask)
		if err != nil {
			views.Error(err.Error(), http.StatusBadRequest, w)
			return
		}

		isTaskRunning, err := sv.IsTaskRunning(r.Context(), existingTask.ID)
		if err != nil {
			views.Error("task does not exist", http.StatusBadRequest, w)
			return
		}

		if !isTaskRunning {
			views.Error("task is not running", http.StatusBadRequest, w)
			return
		}

		cancelled, err := sv.CancelTask(r.Context(), wg, existingTask.ID)
		if err != nil || !cancelled {
			views.Error("task didn't get cancelled", http.StatusInternalServerError, w)
			return
		}

		views.Success("task cancelled successfully", http.StatusOK, w)
	})
}

func ping(sv task.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			views.Error(pkg.ErrMethodNotAllowed.Error(), http.StatusMethodNotAllowed, w)
			return
		}
		views.Success("pong", http.StatusCreated, w)
	})
}

func GenerateTaskHandler(r *http.ServeMux, wg *sync.WaitGroup, sv task.Service) {
	r.Handle("/api/v1/task/ping", ping(sv))
	r.Handle("/api/v1/task/create", create(sv))
	r.Handle("/api/v1/task/submit", submit(sv, wg))
	r.Handle("/api/v1/task/view", view(sv))
	r.Handle("/api/v1/task/viewall", viewAll(sv))
	r.Handle("/api/v1/task/cancel", cancel(sv, wg))
}

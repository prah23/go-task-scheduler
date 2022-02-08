package task

import (
	"context"
	"math"
	"path/filepath"
	"sync"
	"time"
)

type Service interface {
	CreateTask(ctx context.Context, task *Task) (*Task, error)

	RunTask(command string, ctx context.Context, wg *sync.WaitGroup, taskID string) (*Task, error)

	CancelTask(ctx context.Context, wg *sync.WaitGroup, taskID string) (bool, error)

	GetTask(ctx context.Context, taskID string) (*Task, error)

	GetAllTasks(ctx context.Context) ([]Task, error)

	IsTaskRunning(ctx context.Context, taskID string) (bool, error)
}

type service struct {
	repo Repository
}

func NewServiceInstance(r Repository) Service {
	return &service{
		repo: r,
	}
}

func (s service) CreateTask(ctx context.Context, task *Task) (newTask *Task, err error) {
	newTask, err = s.repo.NewTask(ctx, task)
	if err != nil {
		return nil, err
	}
	return newTask, nil
}

func (s service) RunTask(command string, ctx context.Context, wg *sync.WaitGroup, taskID string) (task *Task, err error) {
	task, err = s.repo.GetTask(ctx, taskID)
	if err != nil {
		return nil, err
	}

	filePath, err := filepath.Abs("./" + "scripts" + "/" + task.ID + ".sh")
	if err != nil {
		return nil, err
	}
	timeoutSeconds := math.Min(float64(task.Timeout), 20.0)

	wg.Add(1)
	go ExecuteCommand(command, filePath, wg, ctx, task, time.Duration(timeoutSeconds))
	return task, nil
}

func (s service) CancelTask(ctx context.Context, wg *sync.WaitGroup, taskID string) (cancelled bool, err error) {
	task, err := s.repo.GetTask(ctx, taskID)
	if err != nil {
		return false, err
	}
	err = task.Command.Process.Kill()
	if err != nil {
		return false, err
	}
	task.IsRunning = false
	return true, nil
}

func (s service) GetTask(ctx context.Context, taskID string) (task *Task, err error) {
	task, err = s.repo.GetTask(ctx, taskID)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s service) GetAllTasks(ctx context.Context) (tasks []Task, err error) {
	tasks, err = s.repo.GetAllTasks(ctx)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s service) IsTaskRunning(ctx context.Context, taskID string) (running bool, err error) {
	running, err = s.repo.TaskStatus(ctx, taskID)
	if err != nil {
		return false, err
	}
	return running, nil
}

package task

import (
	"context"
	"errors"
)

type repository struct {
	store *[]*Task
}

func NewRepositoryInstance(st *[]*Task) Repository {
	return &repository{
		store: st,
	}
}

func (r repository) NewTask(ctx context.Context, task *Task) (t *Task, err error) {
	*r.store = append(*r.store, task)
	return task, nil
}

func (r repository) TaskStatus(ctx context.Context, taskID string) (running bool, err error) {
	for _, task := range *r.store {
		if taskID == task.ID {
			return task.IsRunning, nil
		}
	}
	err = errors.New("task does not exist")
	return false, err
}

func (r repository) GetTask(ctx context.Context, taskID string) (task *Task, err error) {
	for _, task := range *r.store {
		if taskID == task.ID {
			return task, nil
		}
	}
	err = errors.New("task does not exist")
	return nil, err
}

func (r repository) GetAllTasks(ctx context.Context) (tasks []Task, err error) {
	var arr []Task
	for _, currentTask := range *r.store {
		arr = append(arr, *currentTask)
	}
	return arr, nil
}

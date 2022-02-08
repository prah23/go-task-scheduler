package task

import "context"

type Repository interface {
	NewTask(ctx context.Context, task *Task) (*Task, error)

	TaskStatus(ctx context.Context, taskID string) (bool, error)

	GetTask(ctx context.Context, taskID string) (*Task, error)

	GetAllTasks(ctx context.Context) ([]Task, error)
}

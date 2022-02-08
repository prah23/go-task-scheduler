package task

import "os/exec"

type Task struct {
	ID          string    `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Timeout     uint      `json:"timeout,omitempty"`
	IsRunning   bool      `json:"isrunning"`
	Output      string    `json:"output,omitempty"`
	Command     *exec.Cmd `json:"-"`
}

type TaskCreateRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Timeout     uint   `json:"timeout,omitempty"`
}

type TaskIDRequest struct {
	ID string `json:"id"`
}

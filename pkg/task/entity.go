package task

type Task struct {
	ID          uint   `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Timeout     uint   `json:"timeout,omitempty"`
	IsRunning   bool   `json:"isrunning,omitempty"`
	Output      string `json:"output,omitempty"`
}

package task

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

func ExecuteCommand(command, filePath string, wg *sync.WaitGroup, ctx context.Context, t *Task, timeout time.Duration) {
	defer wg.Done()

	ct, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	t.Command = exec.CommandContext(ct, command, filePath)

	t.IsRunning = true

	out, err := t.Command.CombinedOutput()

	t.IsRunning = false

	if ct.Err() == context.DeadlineExceeded {
		t.Output = "Task timed out"
	}

	if err != nil {
		t.Output = err.Error()
		return
	}

	t.Output = string(out)
}

func StoreScriptInFile(contents []byte, taskID string) (length int, err error) {
	filePath, err := filepath.Abs("./" + "scripts" + "/" + taskID + ".sh")

	file, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	length, err = file.WriteString(string(contents))
	if err != nil {
		return 0, err
	}

	return length, nil
}

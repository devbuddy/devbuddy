package helpers

import (
	"fmt"
	"os"
	"path"
	"time"
)

type DebugLogWriter struct {
	path        string
	maxFileSize int64
}

func NewDebugLogWriter() *DebugLogWriter {
	filename := fmt.Sprintf("devbuddy-debug-%d.log", os.Getuid())

	return &DebugLogWriter{
		path:        path.Join(os.TempDir(), filename),
		maxFileSize: 1 * 1024 * 1024, // 1MB
	}
}

func (w *DebugLogWriter) Write(buffer []byte) {
	if len(buffer) == 0 {
		return
	}

	w.rotate()

	err := w.append(buffer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Debug log error: %s", err)
	}
}

func (w *DebugLogWriter) rotate() {
	fi, err := os.Stat(w.path)
	if err != nil {
		return
	}

	if fi.Size() > w.maxFileSize {
		_ = os.Rename(w.path, w.path+".old")
	}
}

func (w *DebugLogWriter) append(buffer []byte) error {
	file, err := os.OpenFile(w.path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return err
	}

	defer func() {
		_ = file.Close()
	}()

	_, err = file.WriteString(fmt.Sprintf("\n------ %s ------\n", time.Now()))
	if err != nil {
		return err
	}

	_, err = file.Write(buffer)
	if err != nil {
		return err
	}

	return nil
}

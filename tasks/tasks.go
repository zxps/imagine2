package tasks

import (
	"imagine2/models"
	"time"
)

// Task ...
type Task struct {
	Name    string
	Started int64
	Stop    bool
}

// Manager ...
type Manager struct {
	Running       map[string]Task
	DeletedFileID chan int
	CreatedFileID chan int
	RenderedPath  chan string
}

// AddRunning ...
func (m *Manager) AddRunning(name string) {
	m.Running[name] = Task{
		Name:    name,
		Started: time.Now().Unix(),
		Stop:    false,
	}
}

// Stop ...
func (m *Manager) Stop(name string) {

}

// IsRunning ...
func (m *Manager) IsRunning(name string) bool {
	if _, ok := m.Running[name]; ok {
		return true
	}

	return false
}

// TaskManager ...
var TaskManager Manager = Manager{
	Running:       make(map[string]Task),
	DeletedFileID: make(chan int, 100),
	CreatedFileID: make(chan int, 100),
	RenderedPath:  make(chan string, 100),
}

// NotifyFileDelete ...
func NotifyFileDelete(file models.File) {
	TaskManager.DeletedFileID <- file.ID
}

// NotifyFileCreated ...
func NotifyFileCreated(file models.File) {
	TaskManager.CreatedFileID <- file.ID
}

// NotifyPathRendered ...
func NotifyPathRendered(path string) {
	TaskManager.RenderedPath <- path
}

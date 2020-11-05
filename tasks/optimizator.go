package tasks

import (
	"imagine2/files"
	"runtime"
)

const taskOptimizator = "optimizator"

// StartOptimizator ...
func StartOptimizator() {
	if TaskManager.IsRunning(taskOptimizator) {
		return
	}

	TaskManager.AddRunning(taskOptimizator)

	go onFilesDelete(TaskManager.DeletedFileID)
	//go onFilesCreate(TaskManager.CreatedFileID)
	//go onFilesRender(TaskManager.RenderedPath)
}

func onFilesRender(renderedPath <-chan string) {
	for {
		_ = <-renderedPath

		runtime.GC()
	}
}

func onFilesCreate(createdFileID <-chan int) {
	for {
		_ = <-createdFileID
		runtime.GC()
	}
}

func onFilesDelete(deletedFileID <-chan int) {
	for {
		_ = <-deletedFileID
		files.RemoveEmptyDirectories(files.GetFilespath())
		files.RemoveEmptyDirectories(files.GetCachePath())
		runtime.GC()
	}
}

package watcher

import (
	"log"
)

// Send operations by rest API
type logOperationManager struct{}

func NewLogOperationManager()OperationManager {
	return logOperationManager{}
}

func (rom logOperationManager) CopyFile(path string) {
	log.Println("Copy file",path)
}

func (rom logOperationManager) CopyFolder(path string) {
	log.Println("Copy folder",path)
}

func (rom logOperationManager) Delete(path string){
	log.Println("Delete",path)
}


package watcher

type OperationManager interface{
	CopyFile(path string)
	CopyFolder(path string)
	Delete(path string)
}

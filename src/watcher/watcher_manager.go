package watcher

import (
	"errors"
	"github.com/jotitan/fsnot-poc/config"

	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type WatchManager struct {
	watcher *fsnotify.Watcher
	// for each watch folder,
	conf config.Config
	// Store file which is currently written.
	checkWriteDone map[string]struct{}
	operationManager OperationManager
}

func (wm * WatchManager)isWriteDone(path string)bool{
	if _,exist := wm.checkWriteDone[path] ; exist {
		delete(wm.checkWriteDone,path)
		return true
	}
	wm.checkWriteDone[path] = struct{}{}
	return false
}

func NewWatchManager(conf config.Config)(*WatchManager,error){
	if watcher,err := fsnotify.NewWatcher() ; err == nil {
		return  &WatchManager{
			conf:conf,
			watcher:watcher,
			checkWriteDone: make(map[string]struct{}),
			operationManager: NewRestOperationManager(conf),
		},nil
	}else{
		return nil,err
	}
}

func (wm *WatchManager)stopWatch(path string) {
	wm.watcher.Remove(path)
}

func (wm *WatchManager)watch(path string){
	wm.watcher.Add(path)
	log.Println("Add watch",path)
	// Watch also sub folder
	if f,err := os.Open(path) ; err == nil {
		defer f.Close()
		if files,err := f.Readdir(-1) ; err == nil {
			for _,file := range files {
				if file.IsDir() {
					filename := filepath.Join(path, file.Name())
					wm.watch(filename)
				}
			}
		}
	}
}

func (wm WatchManager)findBackupFolder(path string)(string,string,error){
	// Search where to backup file/folder based on root path
	for _,folder := range wm.conf.Watchers {
		if strings.HasPrefix(path,folder.WatchedFolder){
			return folder.BackupFolder,folder.WatchedFolder,nil
		}
	}
	return "","",errors.New("bad folder")
}

func (wm WatchManager)Launch(){
	go func(){
		for {
			if value,hasMore := <- wm.watcher.Events ; hasMore {
				log.Println(value)
				switch {
				case isCreate(value.Op):
					wm.create(value.Name,true)
					break
				case isWrite(value.Op):
					wm.create(value.Name,false)
					break
				case isDelete(value.Op):
					wm.delete(value.Name)
					break
				}
			}else{
				break
			}
		}

		wm.watcher.Close()
	}()
	for _,watcher := range wm.conf.Watchers {
		wm.watch(watcher.WatchedFolder)
	}

	log.Println("Backup service launch, wait changes...")
	<- make(chan struct{})
}

func isCreate(op fsnotify.Op)bool{
	return isOperation(op,fsnotify.Create)
}

func isWrite(op fsnotify.Op)bool{
	return isOperation(op,fsnotify.Write)
}

func isDelete(op fsnotify.Op)bool{
	return isOperation(op,fsnotify.Remove)
}

func isOperation(op fsnotify.Op,mask fsnotify.Op)bool{
	return op&mask == mask
}

// includeFolder manage changes on folder only if true
// To folder creation, do something for creation. To file creation, write, do something on write (isCreation == false)
// Write file work in two steps : open file to write inside and close. Only copy when second one appends
func (wm * WatchManager)create(path string, isCreation bool){
	if stat,err := os.Lstat(path) ; err == nil {
		if stat.IsDir() {
			if isCreation {
				wm.watch(path)
				wm.copyFolder(path, true)
			}
		} else {
			if !isCreation && wm.isWriteDone(path){
				wm.copyFile(path)
			}
		}
	}else{
		log.Println("ERROR READ FILE",err)
	}
}

// if launchOperation is true, launch copyFolder operation, other not
func (wm * WatchManager)copyFolder(path string,launchOperation bool){
	if launchOperation {
		wm.operationManager.CopyFolder(path)
	}
	if dir,err := os.Open(path) ; err == nil {
		defer dir.Close()
		files,_ := dir.Readdir(-1)
		for _,file := range files {
			subPath := filepath.Join(path,file.Name())
			if file.IsDir() {
				wm.watch(subPath)
				wm.copyFolder(subPath,false)
			}
		}
	}
}

func (wm * WatchManager)copyFile(path string){
	wm.operationManager.CopyFile(path)
}

func (wm * WatchManager)delete(path string){
	wm.operationManager.Delete(path)
	// Check in backup folder if it's a folder, if yes, remove recursive, otherwise, remove file
	/*if stat,err := os.Lstat(path) ; err == nil {
		if stat.IsDir() {
			wm.watcher.Remove(path)
		}
	}else{
		log.Fatal("Impossible to Delete",path,err)
	}*/
}
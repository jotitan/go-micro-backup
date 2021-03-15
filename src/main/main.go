package main

import (
	"errors"
	"github.com/jotitan/fsnot-poc/config"

	"github.com/fsnotify/fsnotify"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main(){
	conf := config.NewConfig("micro-backup.json")
	watch,err  := NewWatchManager(conf)
	if err != nil {
		log.Fatal("Fail !!!",err)
	}
	watch.launch()
}

type WatchManager struct {
	watcher *fsnotify.Watcher
	// for each watch folder,
	conf config.Config
	// Store file which is currently written.
	checkWriteDone map[string]struct{}
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
		return  &WatchManager{conf:conf,watcher:watcher,checkWriteDone: make(map[string]struct{})},nil
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

func (wm WatchManager)launch(){
	go func(){
		for {
			if value,hasMore := <- wm.watcher.Events ; hasMore {
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
// Write file work in two steps : open file to write inside and close. Only copy when second one appends
func (wm * WatchManager)create(path string, includeFolder bool){
	if stat,err := os.Lstat(path) ; err == nil {
		if backupFolder,watchedFolder,err := wm.findBackupFolder(path) ; err == nil {
			if stat.IsDir() {
				if  includeFolder {
					wm.watch(path)
					wm.copyFolder(path, backupFolder, watchedFolder)
				}
			} else {
				if !includeFolder && wm.isWriteDone(path){
					wm.copyFile(path, backupFolder, watchedFolder)
				}
			}
		}
	}else{
		log.Println("ERROR READ FILE",err)
	}
}

func (wm * WatchManager)copyFolder(path,backupFolder,watchedFolder string){
	folderName := filepath.Join(backupFolder,strings.Replace(path,watchedFolder,"",-1))
	if err := os.MkdirAll(folderName,os.ModePerm) ; err != nil {
		log.Fatal("Impossible to create folder",folderName,err)
		return
	}
	log.Println("CREATE FOLDER",folderName)

	// Browse files inside
	if dir,err := os.Open(path) ; err == nil {
		files,_ := dir.Readdir(-1)
		for _,file := range files {
			subPath := filepath.Join(path,file.Name())
			if file.IsDir() {
				wm.watch(subPath)
				wm.copyFolder(subPath,backupFolder,watchedFolder)
			}else{
				wm.copyFile(subPath,backupFolder,watchedFolder)
			}
		}
	}
}

func (wm * WatchManager)copyFile(path,backupFolder,watchedFolder string){
	filename := filepath.Join(backupFolder,strings.Replace(path,watchedFolder,"",-1))
	if out,err := os.OpenFile(filename,os.O_CREATE|os.O_RDWR|os.O_TRUNC,os.ModePerm) ; err == nil {
		defer out.Close()
		if in,err := os.Open(path) ; err == nil {
			defer in.Close()
			io.Copy(out,in)
			log.Println("COPY FILE from",path,"to",filename)
		}
	}
}

func (wm * WatchManager)delete(path string){
	// Check in backup folder if it's a folder, if yes, remove recursive, otherwise, remove file
	if backupFolder,watchedFolder,err := wm.findBackupFolder(path) ; err == nil {
		filename := filepath.Join(backupFolder,strings.Replace(path,watchedFolder,"",-1))

		if stat,err := os.Lstat(filename) ; err == nil {
			if stat.IsDir() {
				if err := os.RemoveAll(filename) ; err == nil {
					log.Println("DELETE FOLDER", path)
				}else{
					log.Fatal("Impossible to delete folder",path,err)
				}
			}else{
				if err := os.Remove(filename) ; err == nil {
					log.Println("DELETE FILE",path)
				}else{
					log.Fatal("Impossible to delete file",path,err)
				}
			}
		}else{
			log.Fatal("Impossible to delete",path,err)
		}
	}


}
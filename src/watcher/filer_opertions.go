package watcher

import (
	"errors"
	"github.com/jotitan/fsnot-poc/config"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Send operations by rest API
type filerOperationManager struct{
	watchersConf config.Config
	backupFolder string
}

func NewFilerOperationManager(conf config.Config)OperationManager {
	return filerOperationManager{watchersConf:conf}
}

func (rom filerOperationManager)findBackupFolder(path string)(string,string,error){
	// Search where to backup file/folder based on root path
	for _,folder := range rom.watchersConf.Watchers {
		if strings.HasPrefix(path,folder.WatchedFolder){
			return folder.BackupFolder,folder.WatchedFolder,nil
		}
	}
	return "","",errors.New("bad folder")
}

func (rom filerOperationManager) CopyFolder(path string) {
	backupFolder,sourceFolder,_ := rom.findBackupFolder(path)
	filename := filepath.Join(backupFolder,strings.Replace(path,sourceFolder,"",-1))

	if err := os.MkdirAll(filename,os.ModePerm) ; err != nil {
		log.Fatal("Impossible to create folder",filename,err)
		return
	}
	log.Println("CREATE FOLDER",filename)

	// Browse files inside
	if dir,err := os.Open(path) ; err == nil {
		files,_ := dir.Readdir(-1)
		for _,file := range files {
			subPath := filepath.Join(path,file.Name())
			if file.IsDir() {
				rom.CopyFolder(subPath)
			}else{
				rom.CopyFile(subPath)
			}
		}
	}
}

// The path is clean and relative
func (rom filerOperationManager) CopyFile(path string){
	backupFolder,sourceFolder,_ := rom.findBackupFolder(path)
	filename := filepath.Join(backupFolder,strings.Replace(path,sourceFolder,"",-1))
	if out,err := os.OpenFile(filename,os.O_CREATE|os.O_RDWR|os.O_TRUNC,os.ModePerm) ; err == nil {
		defer out.Close()
		if in,err := os.Open(path) ; err == nil {
			defer in.Close()
			io.Copy(out,in)
			log.Println("COPY FILE from",path,"to",filename)
		}
	}
}

// The path is clean and relative
func (rom filerOperationManager) Delete(path string){
	filename := filepath.Join(rom.backupFolder,path)

	if stat,err := os.Lstat(filename) ; err == nil {
		if stat.IsDir() {
			if err := os.RemoveAll(filename) ; err == nil {
				log.Println("DELETE FOLDER", path)
			}else{
				log.Fatal("Impossible to Delete folder",path,err)
			}
		}else{
			if err := os.Remove(filename) ; err == nil {
				log.Println("DELETE FILE",path)
			}else{
				log.Fatal("Impossible to Delete file",path,err)
			}
		}
	}else{
		log.Fatal("Impossible to Delete",path,err)
	}
}


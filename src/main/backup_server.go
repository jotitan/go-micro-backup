package main

import (
	"github.com/jotitan/fsnot-poc/config"
	"github.com/jotitan/fsnot-poc/watcher"
	"log"
	"net/http"
	"strings"
)

func main(){
	conf := config.NewConfig("micro-backup.json")
	server := newBackupServer(conf)
	server.launch()
}

type backupServer struct {
	operationManager watcher.OperationManager
}

func newBackupServer(conf config.Config) backupServer {
	return backupServer{watcher.NewFilerOperationManager(conf)}
}

func (bs backupServer)launch(){
	server := http.NewServeMux()
	server.HandleFunc("/file",bs.copyFile)
	server.HandleFunc("/folder",bs.copyFolder)
	server.HandleFunc("/health",bs.health)
	server.HandleFunc("/",bs.delete)

	log.Println("Launch server on port 9002")
	http.ListenAndServe(":9002",server)
}

func (bs backupServer)copyFile(w http.ResponseWriter,r * http.Request){
	if r.Method != http.MethodPost {
		http.Error(w,"Only POST allowed",405)
		return
	}
	if path := r.FormValue("path") ; !strings.EqualFold(path,"") {
		bs.operationManager.CopyFile(path)
	}else{
		http.Error(w,"path required",400)
	}
}

func (bs backupServer)copyFolder(w http.ResponseWriter,r * http.Request){
	if r.Method != http.MethodPost {
		http.Error(w,"Only POST allowed",405)
		return
	}
	if path := r.FormValue("path") ; !strings.EqualFold(path,"") {
		bs.operationManager.CopyFolder(path)
	}else{
		http.Error(w,"path required",400)
	}
}

func (bs backupServer)health(w http.ResponseWriter,r * http.Request){
	w.WriteHeader(200)
	w.Write([]byte("up"))
}

func (bs backupServer)delete(w http.ResponseWriter,r * http.Request){
	if r.Method != http.MethodDelete {
		http.Error(w,"Only DELETE allowed",405)
		return
	}
	if path := r.FormValue("path") ; !strings.EqualFold(path,"") {
		bs.operationManager.Delete(path)
	}else{
		http.Error(w,"path required",400)
	}
}
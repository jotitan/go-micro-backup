package watcher

import (
	"fmt"
	"github.com/jotitan/fsnot-poc/config"
	"log"
	"net/http"
	"net/url"
)

// Send operations by rest API
type restOperationManager struct{
	url string
}

func NewRestOperationManager(conf config.Config)OperationManager {
	return restOperationManager{url:conf.BackupRestUrl}
}

func (rom restOperationManager) CopyFile(path string) {
	log.Println("Rest Copy File",path)
	restUrl := fmt.Sprintf("%s/file",rom.url)
	values := url.Values{"path":{path}}
	if _,err := http.DefaultClient.PostForm(restUrl,values) ; err != nil {
		log.Println("Error",err)
	}
}

func (rom restOperationManager) CopyFolder(path string) {
	log.Println("Rest Copy Folder",path)
	restUrl := fmt.Sprintf("%s/folder",rom.url)
	values := url.Values{"path":{path}}
	if _,err := http.DefaultClient.PostForm(restUrl,values) ; err != nil {
		log.Println("Error",err)
	}
}

func (rom restOperationManager) Delete(path string){
	log.Println("Delete",path)
	restUrl := fmt.Sprintf("%s/path=%s",rom.url,path)
	r,_ := http.NewRequest(http.MethodDelete,restUrl,nil)
	if _,err := http.DefaultClient.Do(r); err != nil {
		log.Println("Error",err)
	}
}


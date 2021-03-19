package watcher

import (
	"fmt"
	"github.com/jotitan/fsnot-poc/config"
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
	restUrl := fmt.Sprintf("%s/file",rom.url)
	values := url.Values{"path":{path}}
	http.DefaultClient.PostForm(restUrl,values)
}

func (rom restOperationManager) CopyFolder(path string) {
	restUrl := fmt.Sprintf("%s/folder",rom.url)
	values := url.Values{"path":{path}}
	http.DefaultClient.PostForm(restUrl,values)
}

func (rom restOperationManager) Delete(path string){
	restUrl := fmt.Sprintf("%s/path=%s",rom.url,path)
	r,_ := http.NewRequest(http.MethodDelete,restUrl,nil)
	http.DefaultClient.Do(r)
}


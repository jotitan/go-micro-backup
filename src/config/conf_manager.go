package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	Watchers []struct{
		WatchedFolder string `json:"watched_folder"`
		BackupFolder string `json:"backup_folder"`
	} `json:"watchers"`
}

func NewConfig(path string) Config{
	c := Config{}
	if data,err := ioutil.ReadFile(path) ; err == nil {
		json.Unmarshal(data,&c)
		log.Println("Read config",len(c.Watchers),"path(s)")
	}else{
		log.Println("Error",err)
	}
	return c
}
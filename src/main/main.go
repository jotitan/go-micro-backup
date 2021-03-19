package main

import (
	"github.com/jotitan/fsnot-poc/config"
	"github.com/jotitan/fsnot-poc/watcher"

	"log"
)

func main(){
	conf := config.NewConfig("micro-backup.json")
	watch,err  := watcher.NewWatchManager(conf)
	if err != nil {
		log.Fatal("Fail !!!",err)
	}
	watch.Launch()
}
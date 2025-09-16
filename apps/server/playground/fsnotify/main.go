package main

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		log.Println("Watching")
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				log.Println("event:", event, event.Op.String())

				if event.Has(fsnotify.Write) {
					log.Println("modified file:", event.Name)
				}
				if event.Has(fsnotify.Rename) {
					log.Println("renamed file:", event.Name)
				}
				if event.Has(fsnotify.Create) {
					log.Println("created file:", event.Name)
					//log.Println(event.renamedFrom)
				}

				if event.Has(fsnotify.Rename) && event.Has(fsnotify.Create) {
					log.Println("Eureka")
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("./tmp")
	if err != nil {
		log.Fatal(err)
	}

	<-make(chan struct{})
}

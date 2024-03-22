package Watchdogs

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"time"
)

func FileWatchdog(filePath string, callback callback, timeout time.Duration) (destructor, error) {
	done := make(chan bool)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return func() {
			done <- true
		}, err
	}
	defer func(watcher *fsnotify.Watcher) {
		err := watcher.Close()
		if err != nil {
			panic(err)
		}
	}(watcher)

	err = watcher.Add(filePath)
	if err != nil {
		return func() {
			done <- true
		}, err
	}

	ticker := time.NewTicker(timeout)

	go func() {
		for done != nil {
			select {
			case <-ticker.C:
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}
					if event.Op&fsnotify.Write == fsnotify.Write {

						fileData, err := os.ReadFile(filePath)
						if err != nil {
							panic(err)
						}

						stringFileData := string(fileData)

						err = callback(&stringFileData)
					}
				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}
					log.Println("error:", err)
				}
			}
		}
	}()

	return func() {
		done <- true
	}, nil
}

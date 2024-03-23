package Watchdogs

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"github.com/fsnotify/fsnotify"
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

	go func() {
		var savedHash []byte

		ticker := time.NewTicker(timeout)
		for done != nil {
			println("Checking file")
			select {
			case <-ticker.C:
			}

			// calculate sha1 hash of file
			fileData, err := os.ReadFile(filePath)
			if err != nil {
				panic(err)
			}

			hashCreator := sha1.New()
			_, err = hashCreator.Write(fileData)
			if err != nil {
				panic(err)
			}

			hash := hashCreator.Sum(nil)

			if bytes.Equal(hash, savedHash) {
				continue
			}
			println("File changed")
			fmt.Printf("1: %x ", hash)
			fmt.Printf("2: %x\n", savedHash)
			savedHash = hash

			stringFileData := string(fileData)

			err = callback(&stringFileData)
			if err != nil {
				panic(err)
			}
		}

	}()

	return func() {
		done <- true
	}, nil
}

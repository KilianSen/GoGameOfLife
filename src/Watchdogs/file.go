package Watchdogs

import (
	"bytes"
	"crypto/sha1"
	"os"
	"time"
)

func getHash(data []byte) []byte {
	hashCreator := sha1.New()
	_, err := hashCreator.Write(data)
	if err != nil {
		panic(err)
	}

	return hashCreator.Sum(nil)
}

func FileWatchdog(filePath string, callback callback, timeout time.Duration) (Destructor, error) {
	done := make(chan bool)

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var savedHash []byte = getHash(fileData)

	go func() {

		ticker := time.NewTicker(timeout)
		for done != nil {
			println("Checking file")
			select {
			case <-ticker.C:
			}

			fileData, err := os.ReadFile(filePath)
			if err != nil {
				panic(err)
			}

			hash := getHash(fileData)

			if bytes.Equal(hash, savedHash) {
				continue
			}

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

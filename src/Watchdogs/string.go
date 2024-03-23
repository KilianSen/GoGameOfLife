package Watchdogs

import (
	"bytes"
	"crypto/sha1"
	"time"
)

func StringWatchdog(str *string, callback callback, timeout time.Duration) Destructor {
	var storedHash []byte

	ticker := time.NewTicker(timeout)
	done := make(chan bool)

	go func() {
		defer ticker.Stop()
		for done != nil {
			select {
			case <-ticker.C:
				hashCreator := sha1.New()
				_, err := hashCreator.Write([]byte(*str))
				if err != nil {
					panic(err)
				}

				hash := hashCreator.Sum(nil)

				if bytes.Equal(hash, storedHash) {
					continue
				}

				storedHash = hash

				err = callback(str)
				if err != nil {
					panic(err)
				}
			}
		}
	}()

	return func() {
		ticker.Stop()
		<-done
	}
}

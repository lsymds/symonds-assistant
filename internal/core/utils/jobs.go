package utils

import (
	"log"
	"time"
)

// RunInBackground runs a function every duration specified.
func RunInBackground(name string, f func() error, every time.Duration) {
	go func() {
		for {
			// this is annoying, but the only way to recover and carry on
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("background process(%v): %v", name, r)
					}
				}()

				log.Printf("background process(%v): starting", name)

				if err := f(); err != nil {
					log.Printf("background process(%v): err: %v", name, err.Error())
				}

				log.Printf("background process(%v): ending", name)

				<-time.After(every)
			}()
		}
	}()
}

package main

import (
	"log"
	"time"

	"github.com/sevlyar/go-daemon"
)

// IPLine is a struct that stores the IP and current time.
type IPLine struct {
	IP   string
	Time time.Time
}

func main() {
	cntxt := &daemon.Context{}
	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()
	runService()
}

func runService() {
	c := make(chan IPLine)
	quit := make(chan bool)
	go func() {
		for {
			c1 := getIP()
			c <- c1
			// Check the IP every 20 Minute.
			time.Sleep(20 * time.Minute)
		}
	}()
	go func() {
		for {
			select {
			case c1 := <-c:
				// Channel to log the IP Address.
				logToFile(c1)
			}
		}
	}()

	// Waiting channel for running the program indefintely
	if !<-quit {

	}
}

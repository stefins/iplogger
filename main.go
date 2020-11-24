package main

import (
	"time"
)

// IPLine is a struct that stores the IP and current time.
type IpLine struct {
	Ip   string
	Time time.Time
}

func main() {
	c := make(chan IpLine)
	quit := make(chan bool)
	go func() {
		for {
			c1 := getIp()
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

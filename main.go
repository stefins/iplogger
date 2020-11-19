package main

import (
	"time"
)

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
			time.Sleep(20 * time.Minute)
		}
	}()
	go func() {
		for {
			select {
			case c1 := <-c:
				logToFile(c1)
			}
		}
	}()

	if !<-quit {

	}
}

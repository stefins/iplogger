package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func getIp() (Ip IpLine) {
	// Function to get the current ip address of the device
	resp, err := http.Get("https://ifconfig.me")
	if err != nil {
		// this returns 'Connection Not Available' if internet is not available
		Ip = IpLine{"Connection Not Available", time.Now()}
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	Ip = IpLine{string(body), time.Now()}
	return
}

func logToFile(content IpLine) {
	// Function to log Ip Address to the destination file.
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	createFolder(home)
	ifFileDoesntExist(home, content)
	f, err := os.OpenFile(home+"/.iplogger/log.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	lastIPLine := getLastLine()
	if lastIPLine.Ip != content.Ip {
		// Also logs the current time in RFC1123 format seperated by '*'
		if _, err = f.WriteString(content.Ip + "*" + content.Time.Format(time.RFC1123) + "\n"); err != nil {
			panic(err)
		}
	}
}

func readLineByLine() (ipline []IpLine) {
	// Function to read all the IP address and time from the file
	home, err := os.UserHomeDir()
	file, err := os.Open(home + "/.iplogger/log.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		txt := scanner.Text()
		// Spliting by '*'
		dividedStrings := strings.Split(txt, "*")
		parsedTime, err := time.Parse(time.RFC1123, dividedStrings[1])
		if err != nil {
			panic(err)
		}
		ipline = append(ipline, IpLine{Ip: dividedStrings[0], Time: parsedTime})
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return
}

func getLastLine() (ipline IpLine) {
	// Function to return the last line of the log file
	lines := readLineByLine()
	ipline = lines[len(lines)-1]
	return
}

func ifFileDoesntExist(home string, content IpLine) {
	// Function to create the file if it doesn't exist in the destination
	// (also write the first line)
	if !fileExists(home + "/.iplogger/log.txt") {
		f, err := os.OpenFile(home+"/.iplogger/log.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}
		if _, err = f.WriteString(content.Ip + "*" + content.Time.Format(time.RFC1123) + "\n"); err != nil {
			panic(err)
		}
	}
}

func fileExists(name string) bool {
	// Checking if a file exist or not!
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func createFolder(home string) {
	// Create the new directory if it doesn't exist
	_, err := os.Stat(home + "/.iplogger")
	if os.IsNotExist(err) {
		err := os.Mkdir(home+"/.iplogger", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}

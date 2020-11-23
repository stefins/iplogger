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
	resp, err := http.Get("https://ifconfig.me")
	if err != nil {
		Ip = IpLine{"Connection Not Available", time.Now()}
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	Ip = IpLine{string(body), time.Now()}
	return
}

func logToFile(content IpLine) {
	home, err := os.UserHomeDir()
	ifFileDoesntExist(home, content)
	f, err := os.OpenFile(home+"/.iplogger/log.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	lastIPLine := getLastLine()
	if lastIPLine.Ip != content.Ip {
		if _, err = f.WriteString(content.Ip + "*" + content.Time.Format(time.RFC1123) + "\n"); err != nil {
			panic(err)
		}
	}
}

func readLineByLine() (ipline []IpLine) {
	home, err := os.UserHomeDir()
	file, err := os.Open(home + "/.iplogger/log.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		txt := scanner.Text()
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
	lines := readLineByLine()
	ipline = lines[len(lines)-1]
	return
}

func ifFileDoesntExist(home string, content IpLine) {
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
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

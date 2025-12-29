package main

import (
	"fmt"
	"strconv"
	"strings"
)

type LogEntry struct {
	ip string
	path string
	code int
}

func main() {
	var hardcodedData []string = []string {
		"10.0.0.1,/home,200",
		"10.0.0.2,/about,404",
		"10.0.0.1,/contact,200",
		"10.0.0.3,/pricing,500",
		"10.0.0.2,/login,404",
	}

	LogEnteries := make(map[string]int)
	code404:=0

	for _,val :=range(hardcodedData) {
		splitedData:=strings.Split(val,",")
		ip:=splitedData[0]
		path:=splitedData[1]
		code,err:=strconv.Atoi(splitedData[2])
		if err != nil {
			fmt.Println("Error converting code to int")
		}
		if code == 404 {
			code404++
		}
		LogEnteries[ip]++
		logEntry:=LogEntry{ip,path,code}
		fmt.Println(logEntry)
	}
	maxCount:=0
	maxIP:=""
	for ip,count := range(LogEnteries) {
		if count > maxCount {
			maxCount = count
			maxIP = ip
		}
	}
	fmt.Println("Total 404 errors:",code404)
	fmt.Printf("Most used ip %v with count %v\n",maxIP,maxCount)
}

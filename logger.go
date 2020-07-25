package main

import (
    "log"
    "os"
    "fmt"
    "path/filepath"
)

var (
    LogError  *log.Logger
    LogInfo    *log.Logger
    LogWarn  *log.Logger
)

func init() {
    absPath, err := filepath.Abs("./log")
    if err != nil {
	fmt.Println("Error reading log path:", err)
    }

    file, err := os.OpenFile(absPath+"/httpGo.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {        
	fmt.Println("Error opening log path:", err)
	log.Fatal(err)
    }

    LogInfo   = log.New(file, "INFO: ", log.Ldate|log.Ltime) //log.Ldate|log.Ltime|log.Lshortfile
    LogWarn = log.New(file, "WRN: ", log.Ldate|log.Ltime)
    LogError = log.New(file, "ERR.: ", log.Ldate|log.Ltime)
}

func conditionalLogger(condition bool, msg1, msg2 string) {
	if condition {
		LogInfo.Println(msg1)		
	} else {
		LogError.Println(msg2)	
	}
}
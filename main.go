package main

import (
	"calendar/app"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const (
	inputTimeFormat = "2006-01-02"
)

func main() {
	args := os.Args[1:]
	if len(args) < 3 {
		log.Fatalln("Provide fileName, startDate and endDate as arguments")
	}

	fileName := args[0]
	startDayStr := args[1]
	endDayStr := args[2]

	// Assuming the file is not very very big
	fileBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Failed to read from file: %+v", err)
	}

	startDate, err := time.Parse(inputTimeFormat, startDayStr)
	if err != nil {
		log.Fatalf("Invalid start date format: %+v", err)
	}

	endDate, err := time.Parse(inputTimeFormat, endDayStr)
	if err != nil {
		log.Fatalf("Invalid end date format: %+v", err)
	}
	endDate = endDate.Add((24 * time.Hour) - (1 * time.Second))

	app.Run(fileBytes, startDate, endDate, os.Stdout)
}

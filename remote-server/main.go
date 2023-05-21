package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Results channel
var broadcastResults = make(chan []byte, 1000)
var subscribers = make([]chan []byte, 0)

// Broadcast results to all subscribers
func startBroadcastResults() {
	for result := range broadcastResults {
		for _, subscriber := range subscribers {
			subscriber <- result
		}
	}
}

func subscribeResults(ch chan []byte) {
	subscribers = append(subscribers, ch)
}

func unsubscribeResults(ch chan []byte) {
	newSubscribers := make([]chan []byte, 0)
	for _, subscriber := range subscribers {
		if subscriber == ch {
			continue
		}
		newSubscribers = append(newSubscribers, subscriber)
	}
	subscribers = newSubscribers
}

func startServer(listen string) error {
	// Start httprouter server
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/program", apiProgramSave)
	e.GET("/program", apiProgramGet)
	e.GET("/program/id", apiProgramGetID)
	e.GET("/results", apiResultsStream)
	e.POST("/results", apiResultsPublish)

	e.Static("/", "ui/build")

	go startBroadcastResults()

	return e.Start(listen)
}

func main() {
	log.Println("starting server...")

	// Start http server
	if err := startServer(":8080"); err != nil {
		log.Fatal(err)
	}
}

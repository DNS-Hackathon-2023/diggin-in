package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Results channel
var broadcastResults = make(chan []byte, 1000)

func startServer(listen string) error {
	// Start httprouter server
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/program", apiProgramSave)
	e.GET("/program", apiProgramGet)
	e.GET("/program/id", apiProgramGetID)
	e.GET("/results", apiResultsStream)
	e.POST("/results", apiResultsPublish)
	return e.Start(listen)
}

func main() {
	log.Println("starting server...")

	// Start http server
	if err := startServer(":8080"); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"io/ioutil"
	"log"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

func apiProgramGet(c echo.Context) error {
	p := DefaultProgram()
	program, err := p.Load()
	if err != nil {
		return err
	}
	return c.String(200, program)
}

func apiProgramSave(c echo.Context) error {
	p := DefaultProgram()
	// Read the body into a string
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	if err := p.Save(body); err != nil {
		return err
	}

	if err != nil {
		return err
	}
	return c.String(200, "OK")
}

// Get the sha256 sum of the program
func apiProgramGetID(c echo.Context) error {
	p := DefaultProgram()
	id, err := p.ID()
	if err != nil {
		return err
	}
	return c.String(200, id)
}

// Publish results
func apiResultsPublish(c echo.Context) error {
	// Read body text
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	// Publish to the results channel
	broadcastResults <- body
	return c.String(200, "OK")
}

// Create a websocket handler and streams api results
func apiResultsStream(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		results := make(chan []byte)
		subscribeResults(results)
		defer unsubscribeResults(results)

		for body := range results {
			err := websocket.Message.Send(ws, string(body))
			if err != nil {
				log.Println(err)
				return
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

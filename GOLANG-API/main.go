package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Event struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	UserId int    `json:"user_id"`
}

var events = []Event{}

func main() {

	router := gin.Default()
	router.GET("/events", getEvents)
	router.POST("/events", createEvent)
	router.Run("localhost:8090")
}

func getEvents(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, events)
}

func createEvent(context *gin.Context) {
	var event Event

	if err := context.BindJSON(&event); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"Message": "Can't create an Event"})
		return
	}
	context.IndentedJSON(http.StatusCreated, gin.H{"Message": "Event Created", "Event": event})
	events = append(events, event)
}

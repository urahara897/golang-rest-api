package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"learn-golang/rest-api/models"

	"github.com/gin-gonic/gin"
)

func getEvent(context *gin.Context) {
	eventID, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		fmt.Print(err)
		context.JSON(http.StatusBadRequest, "Could not parse eventID.")
		return
	}

	event, err := models.GetEventByID(eventID)

	if err != nil {
		fmt.Print(err)
		context.JSON(http.StatusInternalServerError, "Could not fetch event.")
		return
	}

	context.JSON(http.StatusOK, event)
}

func getEvents(context *gin.Context) {
	events, err := models.GetAllEvents()

	if err != nil {
		fmt.Print(err)
		context.JSON(http.StatusInternalServerError, "Could not fetch events. Try again later.")
		return
	}
	context.JSON(http.StatusOK, events)
}

func createEvent(context *gin.Context) {
	var event models.Event
	err := context.ShouldBindJSON(&event)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse Request Data.", "error": err})
		return
	}

	event.ID = 1

	err = event.Save()
	if err != nil {
		fmt.Print(err)
		context.JSON(http.StatusInternalServerError, "Could not create event. Try again later.")
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Event Created!", "event": event})
}

func updateEvent(context *gin.Context) {
	eventID, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		fmt.Print(err)
		context.JSON(http.StatusBadRequest, "Could not parse eventID.")
		return
	}

	_, err = models.GetEventByID(eventID)

	if err != nil {
		fmt.Print(err)
		context.JSON(http.StatusBadRequest, "Could not fetch the event.")
		return
	}

	var updatedEvent models.Event
	err = context.ShouldBindJSON(&updatedEvent)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse Request Data.", "error": err})
		return
	}

	updatedEvent.ID = eventID
	err = updatedEvent.Update()

	if err != nil {
		fmt.Print(err)
		context.JSON(http.StatusInternalServerError, "Could not update event. Try again later.")
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Event Updated Successfully"})
}

func deleteEvent(context *gin.Context) {
	eventID, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		fmt.Print(err)
		context.JSON(http.StatusBadRequest, "Could not parse eventID.")
		return
	}

	event, err := models.GetEventByID(eventID)

	if err != nil {
		fmt.Print(err)
		context.JSON(http.StatusBadRequest, "Could not fetch the event.")
		return
	}

	err = event.Delete()

	if err != nil {
		fmt.Print(err)
		context.JSON(http.StatusBadRequest, "Could not delete the event.")
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Event Deleted Successfully"})
}

func deleteAllEvents(context *gin.Context) {
	err := models.DeleteAllEvents()

	if err != nil {
		fmt.Print(err)
		context.JSON(http.StatusBadRequest, "Could not delete all events.")
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Events Deleted Successfully"})
}

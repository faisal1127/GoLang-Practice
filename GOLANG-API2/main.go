package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

var users = []User{}

func main() {
	router := gin.Default()
	router.GET("/users", getUsers)
	router.GET("/users/:id", getEventsByUserId)
	router.POST("/users", createUser)
	router.Run("localhost:8081")
}

func getUsers(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, users)
}

func createUser(context *gin.Context) {
	var user User

	if err := context.BindJSON(&user); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"Message": "Can't create a User"})
		return
	}
	context.IndentedJSON(http.StatusCreated, gin.H{"Message": "User Created", "User": user})
	users = append(users, user)
}

func getEventsByUserId(context *gin.Context) {
	id := context.Param("id")

	resp, err := http.Get("http://localhost:8090/events")
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"Message": "Failed to fetch events"})
		return
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"Message": "Failed to read response"})
		return
	}

	var events []map[string]interface{}

	if err := json.Unmarshal(data, &events); err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"Message": "Failed to parse events"})
		return
	}

	var userEvents []map[string]interface{}
	for _, event := range events {
		if userID, ok := event["user_id"].(float64); ok && int(userID) == atoi(id) {
			userEvents = append(userEvents, event)
		}
	}

	context.IndentedJSON(http.StatusOK, gin.H{"User ID": id, "Events": userEvents})
}

// Helper function to convert string to int
func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

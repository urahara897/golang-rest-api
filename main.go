package main

import (
	"github.com/gin-gonic/gin"

	"learn-golang/rest-api/db"
	"learn-golang/rest-api/routes"
)

func main() {
	db.InitDB()
	server := gin.Default()

	routes.RegisterRoutes(server)

	server.Run(":8080")
}

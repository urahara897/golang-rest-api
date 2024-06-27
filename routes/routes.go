package routes

import (
	"learn-golang/rest-api/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	authenticated := server.Group("/")
	authenticatedAdmin := server.Group("/")

	authenticated.Use(middlewares.Authenticate)
	authenticatedAdmin.Use(middlewares.AuthenticateAdmin)

	authenticated.POST("/events", createEvent)
	authenticated.PUT("/events/:id", updateEvent)
	authenticated.DELETE("/events/:id", deleteEvent)
	authenticated.DELETE("/events", deleteAllEvents)
	authenticated.GET("/events/:id", getEvent)
	authenticated.POST("/events/:id/register", registerForEvents)
	authenticated.DELETE("/events/:id/register", cancelRegistration)

	authenticatedAdmin.GET("/users", getUsers)
	authenticatedAdmin.DELETE("/users/:id", deleteUser)
	authenticatedAdmin.DELETE("/users", deleteAllUsers)

	server.GET("/events", getEvents)
	server.POST("/signup", signup)
	server.POST("/login", login)
	server.POST("/admin/signup", adminSignup)
	server.POST("/admin/login", adminLogin)
}

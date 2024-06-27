package routes

import (
	"fmt"

	"learn-golang/rest-api/utils"

	"github.com/gin-gonic/gin"

	"net/http"

	"learn-golang/rest-api/models"

	"github.com/google/uuid"
)

func checkAdmin(context *gin.Context, admin models.Admin) {
	id := context.GetString("id")
	isAdmin, err := admin.AdminCheck(id)

	if err != nil {
		fmt.Print(err)
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch admin details. Try again later."})
		return
	}

	if !isAdmin {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not an Administrator."})
		return
	}
}

func adminSignup(context *gin.Context) {
	var admin models.Admin

	err := context.ShouldBindJSON(&admin)

	if err != nil {
		fmt.Print(err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse Request Data.", "error": err})
		return
	}

	newUUID := uuid.New().String()
	admin.ID = newUUID

	err = admin.Save()

	if err != nil {
		fmt.Print(err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not save User.", "error": err})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Admin Created successfully.", "Admin": admin})

}

func adminLogin(context *gin.Context) {
	var admin models.Admin

	err := context.ShouldBindJSON(&admin)

	if err != nil {
		fmt.Print(err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse Request Data.", "error": err})
		return
	}

	id, err := admin.ValidateCredentials()

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": err})
		return
	}

	token, err := utils.GenerateToken(admin.Email, id)

	if err != nil {
		fmt.Print(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate token.", "error": err})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Login Successful!!", "token": token})
}

func getUsers(context *gin.Context) {
	var admin models.Admin

	checkAdmin(context, admin)

	events, err := models.GetAllUsers()

	if err != nil {
		fmt.Print(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch users. Try again later."})
		return
	}

	context.JSON(http.StatusOK, events)
}

func deleteUser(context *gin.Context) {
	var admin models.Admin

	userID := context.Param("id")

	checkAdmin(context, admin)

	err := admin.DeleteUser(userID)

	if err != nil {
		fmt.Print(err)
		context.JSON(http.StatusBadRequest, "Could not delete the user.")
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "User Deleted Successfully"})
}

func deleteAllUsers(context *gin.Context) {
	var admin models.Admin

	checkAdmin(context, admin)

	err := models.DeleteAllUsers()

	if err != nil {
		fmt.Print(err)
		context.JSON(http.StatusBadRequest, "Could not delete all users.")
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Users Deleted Successfully"})
}

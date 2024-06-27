package routes

import (
	"fmt"
	"learn-golang/rest-api/models"
	"learn-golang/rest-api/utils"

	"github.com/gin-gonic/gin"

	"net/http"

	"github.com/google/uuid"
)

func signup(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)

	if err != nil {
		fmt.Print(err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse Request Data.", "error": err})
		return
	}

	newUUID := uuid.New().String()
	user.ID = newUUID

	err = user.Save()

	if err != nil {
		fmt.Print(err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not save User.", "error": err})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "User Created successfully.", "User": user})

}

func login(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)

	if err != nil {
		fmt.Print(err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse Request Data.", "error": err})
		return
	}

	id, err := user.ValidateCredentials()

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": err})
		return
	}

	token, err := utils.GenerateToken(user.Email, id)

	if err != nil {
		fmt.Print(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate token.", "error": err})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Login Successful!!", "token": token})
}

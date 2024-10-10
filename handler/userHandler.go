package handler

import (
	"fmt"
	"log"
	"movie_reserve/database"
	"movie_reserve/model"
	"movie_reserve/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserSignUp(context *gin.Context) {
	var newUser model.Users

	if context.Request.Header.Get("Content-Type") != "application/json" {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid Content-Type"})
		fmt.Println("Invalid Content-Type")
		return
	}

	if err := context.ShouldBindJSON(&newUser); err != nil {
		log.Println("Error binding json", err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error binding json"})
		return
	}

	hashedPassword, err := utils.HashPassword(newUser.Password)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error hashing password"})
		return
	}
	newUser.Password = hashedPassword

	_, err = database.DB.Exec("INSERT INTO users (username, password)VALUES($1, $2) ", newUser.Username, newUser.Password)
	if err != nil {
		context.AbortWithStatusJSON(400, "Siging up failed")
		fmt.Println("Error inserting user into database:", err)
	} else {
		context.JSON(http.StatusOK, "User signed up successfully")
	}

}

func Userlogin(context *gin.Context) {
	var newuser model.Users
	err := context.ShouldBindJSON(&newuser)
	if err != nil {
		log.Println("Error binding json", err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error binding json"})
		return
	}
	query := "SELECT password FROM users WHERE username = $1"
	row := database.DB.QueryRow(query, newuser.Username)
	var retrievedpassword string
	err = row.Scan(&retrievedpassword)
	if err != nil {
		log.Fatalf("Error getting retrieved password %v", err)
		return
	}
	passwordValid := utils.ChechHashedPassword(retrievedpassword, newuser.Password)
	if !passwordValid {
		// log.Fatal("Login password is wrong")
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Login password is wrong"})
		return
	}

	token, err := utils.GenerateToken(newuser.Username, newuser.Password)
	if err != nil {
		log.Fatalf("Error generating a token %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate a token."})
	}

	context.JSON(http.StatusOK, gin.H{"message": "Login Successful.", "token": token})
}

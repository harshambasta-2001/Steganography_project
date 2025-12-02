package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harshambasta-2001/Steganography_project/internal"
	"github.com/harshambasta-2001/Steganography_project/utils"
	"github.com/jinzhu/gorm"
)

func (s *APIServer) createuser (c *gin.Context){
	var registeruser internal.RegisterUser

	if err :=c.ShouldBindJSON(&registeruser); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
	}

	hashedPassword,err :=utils.HashPassword(registeruser.Password)
	if err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error": "Failed to hash password"})
	}
	
	user := internal.User{
		Name: registeruser.Name,
		Email:     registeruser.Email,
		Password:  hashedPassword,
	}
	fmt.Print(user)

	if err := s.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "id": user.ID})

}



func (s *APIServer) loginuser (c *gin.Context){
	var loginpayload internal.LoginUser

	if err := c.ShouldBindJSON(&loginpayload); err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})

	}
	email:= loginpayload.Email
	password :=loginpayload.Password

	if email == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}

	user ,err:= s.getUserbyEmail(email)
	if err != nil{
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	} 
		if !utils.CheckPasswordHash(password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := utils.GenerateToken(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}



func (s *APIServer) get_users (c *gin.Context){

	users,err := s.getAllUsers()
	

	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to retrieve all users"})
		return
	}

	c.JSON(http.StatusOK,users)
}




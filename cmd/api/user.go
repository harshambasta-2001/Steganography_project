package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/harshambasta-2001/Steganography_project/internal"
	"github.com/harshambasta-2001/Steganography_project/utils"
	"github.com/jinzhu/gorm"
)

// @Summary Create a new user
// @Description Register a new user with name, email, and password
// @Tags Dashboard
// @Accept  json
// @Produce  json
// @Param   user body internal.RegisterUser true "User Registration Info"
// @Success 201 {object} map[string]interface{} "User created successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /dashboard [post]
func (s *APIServer) createuser(c *gin.Context) {
	var registeruser internal.RegisterUser

	if err := c.ShouldBindJSON(&registeruser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(registeruser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := internal.User{
		Name:     registeruser.Name,
		Email:    registeruser.Email,
		Password: hashedPassword,
	}
	fmt.Print(user)

	if err := s.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "id": user.ID})

}

// @Summary Log in a user
// @Description Authenticate a user and receive a JWT token
// @Tags Dashboard
// @Accept  json
// @Produce  json
// @Param   credentials body internal.LoginUser true "User Login Credentials"
// @Success 200 {object} map[string]string "Login successful"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /dashboard/login [post]
func (s *APIServer) loginuser(c *gin.Context) {
	var loginpayload internal.LoginUser

	if err := c.ShouldBindJSON(&loginpayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	email := loginpayload.Email
	password := loginpayload.Password

	if email == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}

	user, err := s.getUserbyEmail(email)
	if err != nil {
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

// @Summary Get all users
// @Description Retrieve a list of all registered users
// @Tags Dashboard
// @Produce  json
// @Success 200 {array} internal.User
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /dashboard/all-users [get]
func (s *APIServer) get_users(c *gin.Context) {

	users, err := s.getAllUsers()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to retrieve all users"})
		return // Added return
	}

	c.JSON(http.StatusOK, users)
}

// @Summary Delete a user
// @Description Delete a user by their ID
// @Tags Dashboard
// @Produce  json
// @Param   id path int true "User ID"
// @Success 200 {object} map[string]string "User deleted successfully"
// @Failure 400 {object} map[string]string "Invalid ID"
// @Router /dashboard/{id} [delete]
func (s *APIServer) delete_User(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid ID format"})
		return

	}

	_, err = s.delete_user(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User Successfully Deleted"})
}

package main

import (
	// "log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harshambasta-2001/Steganography_project/internal"
	"github.com/harshambasta-2001/Steganography_project/utils"
	"github.com/jinzhu/gorm"
)

func (s *APIServer) createproduct(c *gin.Context) {
	var product internal.PayloadProduct

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId := c.MustGet("userID")

	codes, err := s.get_product_codes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get existing codes"})
		return
	}

	newCode, err := utils.GenerateCode(codes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate unique code"})
		return
	}
	prod := internal.Product{
		Text:   product.Text,
		UserID: userId.(uint),
		Code:   newCode,
	}

	if err := s.db.Create(&prod).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Product created successfully", "Code": prod.Code})

}

func (s *APIServer) extract_text(c *gin.Context) {
	code := c.Param("code")
	userId := c.MustGet("userID")

	product, err := s.get_text_from_code(userId.(uint), code)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve data"})
		return
	}
	// log.Print(product)
	c.JSON(http.StatusOK, gin.H{"message": "Product fetched successfully", "data": product.Text})
}

func (s *APIServer) remove_product(c *gin.Context) {
	code := c.Param("code")
	userId := c.MustGet("userID")

	_, err := s.delete_product(userId.(uint), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

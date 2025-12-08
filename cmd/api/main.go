package main

import (
	"compress/gzip"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/harshambasta-2001/Steganography_project/docs"
	"github.com/harshambasta-2001/Steganography_project/utils"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type APIServer struct {
	addr string
	db   *gorm.DB
}

func NewAPIServer(addr string, db *gorm.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}

}

// @title Steganography Project API
// @version 1.0
// @description This is a server for a steganography application.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email harshambasta12@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8000
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")

	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	apiAddr := os.Getenv("API_ADDR")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}

	}()

	server := NewAPIServer(apiAddr, db)

	router := server.setupRoutes()

	log.Printf("Starting server on %s", apiAddr)
	log.Fatal(router.Run(apiAddr))
}

func (s *APIServer) setupRoutes() *gin.Engine {
	router := gin.Default()
	router.SetTrustedProxies(nil)

	// Add the decompression middleware globally
	router.Use(decompressRequest())

	// Swagger route
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")
	{
		dashboard := v1.Group("/dashboard")
		{
			dashboard.POST("/", s.createuser)
			dashboard.POST("/login", s.loginuser)
			dashboard.GET("/all-users", s.get_users)
			dashboard.DELETE("/:id", s.delete_User)
		}
		product := v1.Group("/product")
		product.Use(utils.AuthMiddleware())
		product.Use(limitRequestBody(1024 * 1024)) // 1 MB limit on DECOMPRESSED data
		{
			product.POST("/", s.createproduct)
			product.GET("/:code", s.extract_text)
			product.DELETE("/:code", s.remove_product)
		}
	}

	return router

}

// limitRequestBody is a middleware that rejects requests with a body larger than maxSize.
func limitRequestBody(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		// The ShouldBindJSON call in the handler will now fail with a specific error
		// if the body is larger than maxSize. We don't need to do anything else here.
		c.Next()
	}
}

// decompressRequest is a middleware that decompresses gzipped request bodies.
func decompressRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("Content-Encoding") == "gzip" {
			reader, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to decompress request body"})
				return
			}
			c.Request.Body = reader
			defer reader.Close()
		}
		c.Next()
	}
}

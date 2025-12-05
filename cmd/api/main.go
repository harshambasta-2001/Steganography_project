package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/harshambasta-2001/Steganography_project/utils"
	"github.com/joho/godotenv"
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
	// router.SetTrustedProxies([]string{"127.0.0.1"})
	router.SetTrustedProxies(nil)

	// router.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "Success connect to database",
	// 	})
	// })

	v1 := router.Group("/api/v1")
	{
		dashboard := v1.Group("/dashboard")
		{
			dashboard.POST("/",s.createuser)
			dashboard.POST("/login",s.loginuser)
			dashboard.GET("/all-users",s.get_users)
			dashboard.DELETE("/:id",s.delete_User)
		}
		product := v1.Group("/product")
		product.Use(utils.AuthMiddleware())
		{
			product.POST("/",s.createproduct)
			product.GET("/:code",s.extract_text)
			product.DELETE("/:code",s.remove_product)
		}
	}

	return router

}

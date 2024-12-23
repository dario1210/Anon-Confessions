package main

import (
	"anon-confessions/cmd/internal/db"
	"anon-confessions/cmd/internal/posts"
	docs "anon-confessions/docs"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Anonymous Confessions API
// @version         1.0
// @description     A privacy-focused backend service that allows users to:
// @description     • Post and manage anonymous confessions.
// @description     • React to posts with likes or dislikes.
// @description     • Leave comments on confessions.
// @description     • Receive real-time updates through WebSocket.
// @description
// @description     The API is designed with RESTful principles, uses SQLite for data storage, and ensures anonymity without storing personal information.

// @termsOfService  http://swagger.io/terms/

// @contact.name    API Support Team
// @contact.url     http://www.swagger.io/support
// @contact.email   support@swagger.io

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @host            localhost:9000
// @BasePath        /api/v1

// @securityDefinitions.basic  BasicAuth

func main() {

	log.Println("Starting the application...")
	log.Println("Initializing database connection...")

	//TODO: READ FROM THE ENVIRONMENT VARIABLE TO GET THE DATABASE NAME AND PASS IT TO THE ConnectToDB FUNCTION
	_, err := db.DbConnection("test.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connected succesfully to the database.")

	router := gin.Default()
	api := router.Group("/api/v1")
	{
		posts.RegisterPostRoutes(api)
	}

	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	log.Println("Starting the HTTP server on :9000...")
	if err := router.Run(":9000"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

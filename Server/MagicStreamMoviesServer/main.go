// marks file as part of main package
package main

// imports packages
import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/gin-gonic/gin"
	"github.com/valentynchystiakow/MagicStreamMovies/Server/MagicStreamMoviesServer/database"
	"github.com/valentynchystiakow/MagicStreamMovies/Server/MagicStreamMoviesServer/routes"
)

// creates main function that runs program
func main() {
	// creates router using gin framework
	router := gin.Default()

	// sets up routes
	router.GET("/hello", func(c *gin.Context) {
		c.String(200, "Hello, MagicStreamMovies!")
	})

	// loads environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: unable to find .env file")
	}

	// defines allowed origins
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")

	// defines origins type
	var origins []string
	// split allowed origins by comma and trim whitespace
	if allowedOrigins != "" {
		origins = strings.Split(allowedOrigins, ",")
		for i := range origins {
			origins[i] = strings.TrimSpace(origins[i])
			log.Println("Allowed Origin:", origins[i])
		}
	} else {
		origins = []string{"http://localhost:5173"}
		log.Println("Allowed Origin: http://localhost:5173")
	}

	// configures cors - allows cross-origin requests to communicate frontend and backend
	config := cors.Config{}
	config.AllowOrigins = origins
	config.AllowMethods = []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"}
	//config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour
	// sets router to use cors
	router.Use(cors.New(config))
	router.Use(gin.Logger())

	// defines client
	var client *mongo.Client = database.Connect()

	// checks if client is connected
	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatalf("Failed to reach server: %v", err)
	}

	// in any case when function ends , closes connection
	defer func() {
		err := client.Disconnect(context.Background())
		if err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}

	}()

	// sets up routes
	routes.SetUpUnprotectedRoutes(router, client)
	routes.SetupProtectedRoutes(router, client)

	// displays error if occurs
	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server", err)
	}

}

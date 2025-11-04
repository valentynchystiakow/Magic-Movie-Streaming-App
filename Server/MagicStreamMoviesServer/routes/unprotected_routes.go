// marks file as part of routes package
package routes

// imports packages
import (
	"github.com/gin-gonic/gin"
	controller "github.com/valentynchystiakow/MagicStreamMovies/Server/MagicStreamMoviesServer/controllers"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// creates function that sets up unprotected routes for unauthenticated users
func SetUpUnprotectedRoutes(router *gin.Engine, client *mongo.Client) {
	// creates route for movies endpoint that handles GET requests to get all movies from database
	router.GET("/movies", controller.GetMovies(client))
	// creates route for register endpoint that handles POST requests to add new user to database
	router.POST("/register", controller.RegisterUser(client))
	// creates route for login endpoint that handles POST requests to login user
	router.POST("/login", controller.LoginUser(client))
	// creates route for logout endpoint that handles POST requests to logout user
	router.POST("/logout", controller.LogoutHandler(client))
	// creates route for genres endpoint that handles GET requests to get all genres
	router.GET("/genres", controller.GetGenres(client))
	// creates route for refresh endpoint that handles POST requests to refresh token
	router.POST("/refresh", controller.RefreshTokenHandler(client))

}

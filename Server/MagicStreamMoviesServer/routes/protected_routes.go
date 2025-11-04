// marks file as part of routes package
package routes

// imports packages
import (
	"github.com/gin-gonic/gin"
	controller "github.com/valentynchystiakow/MagicStreamMovies/Server/MagicStreamMoviesServer/controllers"
	"github.com/valentynchystiakow/MagicStreamMovies/Server/MagicStreamMoviesServer/middleware"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// creates function that sets up protected routes for authenticated users
func SetupProtectedRoutes(router *gin.Engine, client *mongo.Client) {

	// uses AuthMiddleware function to protect routes, requires user to be logged in
	router.Use(middleware.AuthMiddleWare())

	// creates route for movie endpoint that handles GET requests to get certain movie from database
	router.GET("/movie/:movie_id", controller.GetMovie(client))
	// creates route for add-movie endpoint that handles POST requests to add new movie to database
	router.POST("/add-movie", controller.AddMovie(client))
	// creates route for recommended-movies endpoint that handles GET requests to get recommended movies
	router.GET("/recommendedmovies", controller.GetRecommendedMovies(client))
	// creates route for update review endpoint that handles Patch requests to update movie review by imdb id
	router.PATCH("/updatereview/:imdb_id", controller.AdminReviewUpdate(client))
}

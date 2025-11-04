// marks file as part of middleware package
package middleware

// imports packages
import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/valentynchystiakow/MagicStreamMovies/Server/MagicStreamMoviesServer/utils"
)

// creates function that handles authentication middleware
func AuthMiddleWare() gin.HandlerFunc {
	// returns anonymous function that works with gin context
	return func(c *gin.Context) {

		// gets access token
		token, err := utils.GetAccessToken(c)

		// if error occurs, uses JSON to return message with status and aborts request
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// if token is empty, uses JSON to return message with status and aborts request
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}

		// validates token
		claims, err := utils.ValidateToken(token)

		// if error occurs, uses JSON to return message with status and aborts request
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		// sets user id and role to context
		c.Set("userId", claims.UserId)
		c.Set("role", claims.Role)
		// passes request to next middleware handler
		c.Next()

	}
}

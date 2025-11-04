// marks file as part of controllers package
package controllers

// imports packages
import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/valentynchystiakow/MagicStreamMovies/Server/MagicStreamMoviesServer/database"
	"github.com/valentynchystiakow/MagicStreamMovies/Server/MagicStreamMoviesServer/models"
	"github.com/valentynchystiakow/MagicStreamMovies/Server/MagicStreamMoviesServer/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

// creates function that hashes password
func HashPassword(password string) (string, error) {
	// hashes password
	HashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// returns error if occurs
	if err != nil {
		return "", err
	}

	return string(HashPassword), nil

}

// creates function that handles post request to /register-user endpoint
func RegisterUser(client *mongo.Client) gin.HandlerFunc {
	// returns anonymous function that works with gin context - context is used to pass data between handlers
	return func(c *gin.Context) {
		var user models.User
		// uses ShouldBindJSON function to bind json request body to user struct
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
			return
		}
		// defines validator
		validate := validator.New()

		// uses validator to validate user data
		if err := validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}

		// hashes password
		hashedPassword, err := HashPassword(user.Password)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to hash password"})
			return
		}

		// creates special context that cancels request if timeout occurs
		var ctx, cancel = context.WithTimeout(c, 100*time.Second)
		// cancels request when function ends(to prevent memory leaks)
		defer cancel()

		// opens user collection from database
		var userCollection *mongo.Collection = database.OpenCollection("users", client)

		// checks if user already exists
		count, err := userCollection.CountDocuments(ctx, bson.D{{Key: "email", Value: user.Email}})

		// if error occurs uses context to write json response with error message
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing user"})
			return
		}
		// if user already exists uses context to write json response with error message
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return
		}
		// defines user details
		user.UserID = bson.NewObjectID().Hex()
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		user.Password = hashedPassword

		// uses InsertOne function to add new user data to database
		result, err := userCollection.InsertOne(ctx, user)

		// returns error if occurs
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		// uses context to write json response with user data
		c.JSON(http.StatusCreated, result)

	}

}

// creates function that handles post request to /login-user endpoint
func LoginUser(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		// defines userLogin model struct
		var userLogin models.UserLogin

		// uses ShouldBindJSON function to bind json request body to userLogin struct
		if err := c.ShouldBindJSON(&userLogin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalide input data"})
			return
		}

		// when function ends(after 100 seconds) - to prevent memory leaks uses context to cancel request if timeout occurs
		var ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel()

		// opens user collection from database
		var userCollection *mongo.Collection = database.OpenCollection("users", client)

		// checks if user exists in database
		var foundUser models.User
		err := userCollection.FindOne(ctx, bson.D{{Key: "email", Value: userLogin.Email}}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		// compares password with hashed password
		err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(userLogin.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		// generates token and refresh token
		token, refreshToken, err := utils.GenerateAllTokens(foundUser.Email, foundUser.FirstName, foundUser.LastName, foundUser.Role, foundUser.UserID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
			return
		}

		// updates tokens in database
		err = utils.UpdateAllTokens(foundUser.UserID, token, refreshToken, client)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tokens"})
			return
		}

		// sets access_token cookie
		http.SetCookie(c.Writer, &http.Cookie{
			Name:  "access_token",
			Value: token,
			Path:  "/",
			// Domain:   "localhost",
			MaxAge:   86400,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		})
		// sets refresh_token cookie
		http.SetCookie(c.Writer, &http.Cookie{
			Name:  "refresh_token",
			Value: refreshToken,
			Path:  "/",
			// Domain:   "localhost",
			MaxAge:   604800,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		})

		// uses context to write json response with user data
		c.JSON(http.StatusOK, models.UserResponse{
			UserId:    foundUser.UserID,
			FirstName: foundUser.FirstName,
			LastName:  foundUser.LastName,
			Email:     foundUser.Email,
			Role:      foundUser.Role,
			//Token:           token,
			//RefreshToken:    refreshToken,
			FavouriteGenres: foundUser.FavouriteGenres,
		})

	}
}

// creates function that handles post request to /logout-user endpoint
func LogoutHandler(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		// defines userLogout model struct
		var UserLogout struct {
			UserId string `json:"user_id"`
		}

		// uses ShouldBindJSON function to bind json request body to userLogout struct
		err := c.ShouldBindJSON(&UserLogout)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		fmt.Println("User ID from Logout request:", UserLogout.UserId)

		// Clears all tokens for the user
		err = utils.UpdateAllTokens(UserLogout.UserId, "", "", client)

		// Optionally, you can also remove the user session from the database if needed

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error logging out"})
			return
		}
		// c.SetCookie(
		// 	"access_token",
		// 	"",
		// 	-1, // MaxAge negative → delete immediately
		// 	"/",
		// 	"localhost", // Adjust to your domain
		// 	true,        // Use true in production with HTTPS
		// 	true,        // HttpOnly
		// )

		// sets access token in cookie
		http.SetCookie(c.Writer, &http.Cookie{
			Name:  "access_token",
			Value: "",
			Path:  "/",
			// Domain:   "localhost",
			MaxAge:   -1,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		})

		// // Clear the refresh_token cookie
		// c.SetCookie(
		// 	"refresh_token",
		// 	"",
		// 	-1,
		// 	"/",
		// 	"localhost",
		// 	true,
		// 	true,
		// )

		// sets refresh token in cookie
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		})

		c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
	}
}

// creates function that handles post request to /refresh-token endpoint
func RefreshTokenHandler(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// uses special context to cancel request if timeout occurs
		var ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel()

		// gets refresh token from cookie
		refreshToken, err := c.Cookie("refresh_token")

		if err != nil {
			fmt.Println("error", err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to retrieve refresh token from cookie"})
			return
		}

		// defines claim using ValidateRefreshToken function
		claim, err := utils.ValidateRefreshToken(refreshToken)
		if err != nil || claim == nil {
			fmt.Println("error", err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
			return
		}

		// gets user collection from database
		var userCollection *mongo.Collection = database.OpenCollection("users", client)

		// defines filter based on user id
		var user models.User
		err = userCollection.FindOne(ctx, bson.D{{Key: "user_id", Value: claim.UserId}}).Decode(&user)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		// generates new tokens based on user data
		newToken, newRefreshToken, _ := utils.GenerateAllTokens(user.Email, user.FirstName, user.LastName, user.Role, user.UserID)
		err = utils.UpdateAllTokens(user.UserID, newToken, newRefreshToken, client)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating tokens"})
			return
		}

		// sets cookies with new tokens
		c.SetCookie("access_token", newToken, 86400, "/", "localhost", true, true)          // expires in 24 hours
		c.SetCookie("refresh_token", newRefreshToken, 604800, "/", "localhost", true, true) //expires in 1 week

		c.JSON(http.StatusOK, gin.H{"message": "Tokens refreshed"})
	}
}

// marks file as part of utils package
package utils

// imports packages
import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/valentynchystiakow/MagicStreamMovies/Server/MagicStreamMoviesServer/database"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	Role      string
	UserId    string
	jwt.RegisteredClaims
}

// defines secret key
var SECRET_KEY string = os.Getenv("SECRET_KEY")

// defines refresh secret key
var SECRET_REFRESH_KEY string = os.Getenv("SECRET_REFRESH_KEY")

// creates function that generates all tokens
func GenerateAllTokens(email, firstName, lastName, role, userId string) (string, string, error) {
	// generates token
	// defines claims that refer to signed details struct
	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		UserId:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			// registered claims details
			Issuer:    "MagicStream",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	// defines token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// signs token with secret key
	signedToken, err := token.SignedString([]byte(SECRET_KEY))

	// if error occurs, returns empty strings and error
	if err != nil {
		return "", "", err
	}

	// generates refresh token
	// defines claims that refer to signed details struct
	refreshClaims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		UserId:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			// registered claims details
			Issuer:    "MagicStream",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 7 * time.Hour)),
		},
	}

	// defines refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	// signs refresh token with secret key
	signedRefreshToken, err := refreshToken.SignedString([]byte(SECRET_REFRESH_KEY))

	// if error occurs, returns empty strings and error
	if err != nil {
		return "", "", err
	}

	// returns signed token and signed refresh token
	return signedToken, signedRefreshToken, nil

}

// creates function that updates all tokens in database
func UpdateAllTokens(userId, token, refreshToken string, client *mongo.Client) (err error) {
	// creates special context that cancels request if timeout occurs
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	// cancels request when function ends(to prevent memory leaks)
	defer cancel()

	// defines time when token was updated
	updateAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	// defines data to update
	updateData := bson.M{"$set": bson.M{"token": token, "refresh_token": refreshToken, "update_at": updateAt}}

	// opens user collection from database
	var userCollection *mongo.Collection = database.OpenCollection("users", client)

	// updates token data in database
	_, err = userCollection.UpdateOne(ctx, bson.M{"user_id": userId}, updateData)

	// if error occurs, returns error
	if err != nil {
		return err
	}

	return nil
}

// creates function that gets access token
func GetAccessToken(c *gin.Context) (string, error) {
	// authHeader := c.Request.Header.Get("Authorization")
	// if authHeader == "" {
	// 	return "", errors.New("Authorization header is required")
	// }
	// tokenString := authHeader[len("Bearer "):]

	// if tokenString == "" {
	// 	return "", errors.New("Bearer token is required")
	// }

	// gets access token from cookie
	tokenString, err := c.Cookie("access_token")
	if err != nil {

		return "", err
	}

	return tokenString, nil

}

// creates function that validates token
func ValidateToken(tokenString string) (*SignedDetails, error) {

	// defines claims that refer to signed details struct
	claims := &SignedDetails{}

	// defines token using ParseWithClaims function
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// returns secret key as a byte slice type and error
		return []byte(SECRET_KEY), nil
	})
	// if error occurs, returns empty claims and error
	if err != nil {
		return nil, err
	}
	// if token is of invalid type, returns empty claims and error
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, err
	}

	// checks if token is expired
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token has expired")
	}

	return claims, nil
}

// creates function that gets User Id from context
func GetUserIdFromContext(c *gin.Context) (string, error) {
	// gets user id from context
	userId, exists := c.Get("userId")
	// if user id is not found, returns empty string
	if !exists {
		return "", errors.New("user id does not exist in context")
	}
	// if user id is not of string type, returns empty string and error
	id, ok := userId.(string)
	if !ok {
		return "", errors.New("user id does not exists in this context")
	}

	return id, nil
}

// creates function that gets role from context func GetRoleFromContext(c *gin.Context) (string, error) {
func GetRoleFromContext(c *gin.Context) (string, error) {
	// gets role from context
	role, exists := c.Get("role")

	// if role is not found, returns empty string and error
	if !exists {
		return "", errors.New("role does not exists in this context")
	}

	// if role is not of string type, returns empty string and error
	memberRole, ok := role.(string)

	if !ok {
		return "", errors.New("unable to retrieve userId")
	}

	return memberRole, nil

}

// creates function that validates refreshToken
func ValidateRefreshToken(tokenString string) (*SignedDetails, error) {
	// defines claims using SignedDetails struct
	claims := &SignedDetails{}

	// defines token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {

		return []byte(SECRET_REFRESH_KEY), nil
	})

	if err != nil {
		return nil, err
	}

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, err
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("refresh token has expired")
	}

	return claims, nil
}

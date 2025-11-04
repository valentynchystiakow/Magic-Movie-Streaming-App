// marks file as part of controllers package
package controllers

// imports packages
import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms/openai"

	"github.com/valentynchystiakow/MagicStreamMovies/Server/MagicStreamMoviesServer/database"
	"github.com/valentynchystiakow/MagicStreamMovies/Server/MagicStreamMoviesServer/models"
	"github.com/valentynchystiakow/MagicStreamMovies/Server/MagicStreamMoviesServer/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// defines validator
var validate = validator.New()

// creates function that gets movies data(collection) from database, marks function as gin handler in order to be able to handle requests
func GetMovies(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// creates special context that cancels request if timeout occurs - when function ends(after 100 seconds)
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		// uses OpenCollection function from database package to open movies collection from database
		var movieCollection *mongo.Collection = database.OpenCollection("movies", client)

		// gets movies data from database
		var movies []models.Movie
		// creates cursor
		cursor, err := movieCollection.Find(ctx, bson.M{})

		// checks if error occurs
		if err != nil {
			// uses context to write json response with error message
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching movies from database"})
			return
		}
		// closes cursor after function ends(in any case)
		defer cursor.Close(ctx)
		// passes cursor to movies variable if no error occurs
		if err = cursor.All(ctx, &movies); err != nil {
			// uses context to write json response with error message
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode movies from database"})
			return
		}

		// uses context to write json response with movies data
		c.JSON(http.StatusOK, movies)
	}
}

// creates function that gets certain movie data from database, marks function as gin handler in order to be able to handle requests(returns handler function)
func GetMovie(client *mongo.Client) gin.HandlerFunc {
	// returns anonymous function
	return func(c *gin.Context) {
		// creates special context that cancels request if timeout occurs - when function ends(after 100 seconds)
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		// gets movie id from url parameter
		movieID := c.Param("movie_id")
		// if movie id is empty, uses context to write json response with error message
		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID is required"})
			return
		}
		// gets certain movie data from database
		var movie models.Movie
		var movieCollection *mongo.Collection = database.OpenCollection("movies", client)
		// uses FindOne function to get certain movie data by movie id and decodes it
		if err := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID}).Decode(&movie);
		// if error occurs, uses context to write json response with error message
		err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching movie from database"})
			return
		}

		// uses context to write json response with movie data
		c.JSON(http.StatusOK, movie)

	}
}

// creates function that handles post request to /add-movie endpoint to add new movie to a database
func AddMovie(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// creates special context that cancels request if timeout occurs - when function ends(after 100 seconds) - to prevent memory leaks
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		// defines movie variable to store incoming client movie data on server
		var movie models.Movie
		// uses ShouldBindJSON function to bind json request body to movie struct(pointer to movie struct)
		if err := c.ShouldBindJSON(&movie); err != nil {
			// if error occurs, uses context to write json response with error message
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// uses validator to validate movie data
		if err := validate.Struct(movie);
		// if error occurs, uses context to write json response with error message
		err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}

		// defines movieCollection that is type of mongo collection defined in database
		var movieCollection *mongo.Collection = database.OpenCollection("movies", client)

		// uses InsertOne function to add movie data to database
		result, err := movieCollection.InsertOne(ctx, movie)
		// if error occurs, uses context to write json response with error message
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while adding movie to database"})
			return
		}

		// uses context to write json response with movie data
		c.JSON(http.StatusCreated, result)
	}
}

// creates function that updates  movie review as admin
func AdminReviewUpdate(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// uses GetRoleFromContext function to get user role from context
		role, err := utils.GetRoleFromContext(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Role not found in context"})
			return
		}
		// if role is not admin, uses context to write json response with error message
		if role != "ADMIN" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User must be part of the ADMIN role"})
			return
		}
		// defines MovieId from context parameter
		movieId := c.Param("imdb_id")
		// if movie id is empty, uses context to write json response with error message
		if movieId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie Id is required"})
			return
		}

		// defines update review request struct
		var req struct {
			AdminReview string `json:"admin_review"`
		}

		// defines response struct
		var resp struct {
			RankingName string `json:"ranking_name"`
			AdminReview string `json:"admin_review"`
		}

		// uses ShouldBind function to bind json request body to update review request struct
		if err := c.ShouldBind(&req); err != nil {
			// if error occurs, uses context to write json response with error message
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// calls GetReviewRankings function to return sentiment and rank value according to admin review input
		sentiment, rankVal, err := GetReviewRanking(req.AdminReview, client, c)
		if err != nil {
			// uses context to write json response with error message
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while getting review rankings"})
			return
		}

		// defines filter based on movie id
		filter := bson.M{"imdb_id": movieId}

		// defines update that will update fields in database
		update := bson.M{
			"$set": bson.M{
				"admin_review": req.AdminReview,
				"ranking": bson.M{
					"ranking_value": rankVal,
					"ranking_name":  sentiment,
				},
			},
		}

		// uses context to clean resources after performing updates in database
		var ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var movieCollection *mongo.Collection = database.OpenCollection("movies", client)
		// uses UpdateOne function to update movie data in database
		result, err := movieCollection.UpdateOne(ctx, filter, update)

		// if error occurs, uses context to write json response with error message
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while updating movie review in database!"})
			return
		}
		// checks result of update operation, uses context to write json response with error message
		if result.MatchedCount == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Movie not found"})
			return
		}
		// sets response fields with updated data
		resp.RankingName = sentiment
		resp.AdminReview = req.AdminReview

		// uses context to write json response
		c.JSON(http.StatusOK, resp)
	}

}

// creates function that gets review rankings, this function won't work properly if you don't provide your openai api key(you have to deposit 5$ to get started with openai api key)
func GetReviewRanking(admin_review string, client *mongo.Client, c *gin.Context) (string, int, error) {
	// gets all rankings using GetRankings function
	rankings, err := GetRankings(client, c)

	if err != nil {
		return "", 0, err
	}

	// defines sentiment delimited
	sentimentDelimited := ""

	// loops through all rankings to set sentiment delimited
	for _, ranking := range rankings {
		if ranking.RankingValue != 999 {
			sentimentDelimited = sentimentDelimited + ranking.RankingName + ","
		}
	}

	// trims sentiment delimited to remove trailing comma
	sentimentDelimited = strings.Trim(sentimentDelimited, ",")

	// loads env variables
	err = godotenv.Load(".env")

	if err != nil {
		log.Println("Warning: .env file not found")
	}

	// gets openai api key from env
	OpenAiApiKey := os.Getenv("OPENAI_API_KEY")

	// if openai api key is empty, throws error
	if OpenAiApiKey == "" {
		return "", 0, errors.New("could not read OPENAI_API_KEY")
	}

	// defines llm needed for completion using openai
	llm, err := openai.New(openai.WithToken(OpenAiApiKey))

	if err != nil {
		return "", 0, err
	}

	// gets base prompt template from env
	base_prompt_template := os.Getenv("BASE_PROMPT_TEMPLATE")

	// replaces {rankings} with sentiment delimited
	base_prompt := strings.Replace(base_prompt_template, "{rankings}", sentimentDelimited, 1)

	// defines response using llm based on admin review and base prompt
	response, err := llm.Call(context.Background(), base_prompt+admin_review)

	if err != nil {
		return "", 0, err
	}

	// defines rank value
	rankVal := 0

	// loops through all rankings to set rank value
	for _, ranking := range rankings {
		if ranking.RankingName == response {
			rankVal = ranking.RankingValue
			break
		}
	}

	return response, rankVal, nil

}

// creates function that gets rankings
func GetRankings(client *mongo.Client, c *gin.Context) ([]models.Ranking, error) {
	// defines rankings
	var rankings []models.Ranking

	// uses context to cancel request if timeout occurs
	var ctx, cancel = context.WithTimeout(c, 100*time.Second)
	// cancels request when function ends(to prevent memory leaks)
	defer cancel()

	// defines rankingCollection that is type of mongo collection defined in database
	var rankingCollection *mongo.Collection = database.OpenCollection("rankings", client)

	// defines cursor
	cursor, err := rankingCollection.Find(ctx, bson.M{})
	// if error occurs, returns error
	if err != nil {
		return nil, err
	}

	// after function ends, closes cursor
	defer cursor.Close(ctx)

	// passes cursor to rankings variable
	if err := cursor.All(ctx, &rankings); err != nil {
		return nil, err
	}

	return rankings, nil

}

// creates function that handles request to get recommended movies based on user prompt
func GetRecommendedMovies(client *mongo.Client) gin.HandlerFunc {
	// returns anonymous function that works with gin context type
	return func(c *gin.Context) {
		// extracts used id using GetUserIdFromContext function
		userId, err := utils.GetUserIdFromContext(c)
		// if error occurs, uses context to write json response with error message
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error is not found in context"})
			return
		}

		// calls get user favourite genres function to get user favourite genres
		favourite_genres, err := GetUsersFavouriteGenres(userId, client, c)
		// if error occurs, uses context to write json response with error message
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// loads env variables
		err = godotenv.Load(".env")
		// if error occurs, shows warning
		if err != nil {
			log.Println("Warning: .env file not found")
		}
		// defines limit for recommended movies
		var recommendedMovieLimitVal int64 = 5
		// defines limit for recommended movies in string form
		recommendedMovieLimitStr := os.Getenv("RECOMMENDED_MOVIE_LIMIT")

		// checks if limit for recommended movies is set
		if recommendedMovieLimitStr != "" {
			// converts limit for recommended movies to int
			recommendedMovieLimitVal, _ = strconv.ParseInt(recommendedMovieLimitStr, 10, 64)
		}

		// defines find options of movies
		findOptions := options.Find()
		// uses set sort function to set sort order of movies by ranking in descending order
		findOptions.SetSort(bson.D{{Key: "ranking.ranking_value", Value: 1}})

		// sets limit for recommended movies
		findOptions.SetLimit(recommendedMovieLimitVal)

		// defines filter based on favourite genres
		filter := bson.D{
			{Key: "genre.genre_name", Value: bson.D{
				{Key: "$in", Value: favourite_genres},
			}},
		}

		// uses context to cancel request if timeout occurs
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		// cancels request when function ends(to prevent memory leaks)
		defer cancel()

		// defines movieCollection that is type of mongo collection defined in database
		var movieCollection *mongo.Collection = database.OpenCollection("movies", client)

		// defines cursor for iteration
		cursor, err := movieCollection.Find(ctx, filter, findOptions)
		// if error occurs, uses context to write json response with error message
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching movies from database"})
			return
		}

		// after function ends, closes cursor
		defer cursor.Close(ctx)

		// defines recommended movies
		var recommendedMovies []models.Movie
		// passes cursor to recommendedMovies variable
		if err := cursor.All(ctx, &recommendedMovies);
		// if error occurs, uses context to write json response with error message and status
		err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// uses context to write json response with recommended movies and ok status
		c.JSON(http.StatusOK, recommendedMovies)
	}
}

// creates function that gets User Favourite Genres
func GetUsersFavouriteGenres(userId string, client *mongo.Client, c *gin.Context) ([]string, error) {
	// uses context to cancel request if timeout occurs
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	// cancels request when function ends(to prevent memory leaks)
	defer cancel()

	// defines filter based on user id
	filter := bson.D{{Key: "user_id", Value: userId}}

	// defines projection based on favourite genres
	projection := bson.M{"favourite_genres.genre_name": 1,
		"_id": 0}

	// defines options with favourite genres to set projection
	opts := options.FindOne().SetProjection(projection)
	// defines result
	var result bson.M
	// opens user collection from database
	var userCollection *mongo.Collection = database.OpenCollection("users", client)

	// uses FindOne function to get user data by user id and decodes it
	err := userCollection.FindOne(ctx, filter, opts).Decode(&result)
	// if error occurs, returns error
	if err != nil {
		// if error is mongo.ErrNoDocuments, returns empty array
		if err == mongo.ErrNoDocuments {
			return []string{}, nil
		}
	}

	// defines favourite genres array
	favGenresArray, ok := result["favourite_genres"].(bson.A)
	// if nothing is found, returns empty array and error
	if !ok {
		return []string{}, errors.New("unable to retrieve favourite genres")
	}
	// defines genre names array
	var genreNames []string

	// creates loop to iterate through favourite genres array to get genre names and add them to genre names array
	for _, item := range favGenresArray {
		if genreMap, ok := item.(bson.D); ok {
			for _, elem := range genreMap {
				if elem.Key == "genre_name" {
					if name, ok := elem.Value.(string); ok {
						genreNames = append(genreNames, name)
					}
				}
			}
		}
	}

	return genreNames, nil
}

// creates function that gets all genres
func GetGenres(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// uses context to cancel request if timeout occurs
		var ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel()

		// defines genreCollection that is type of mongo collection defined in database
		var genreCollection *mongo.Collection = database.OpenCollection("genres", client)

		// defines cursor to iterate through genres
		cursor, err := genreCollection.Find(ctx, bson.D{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching movie genres"})
			return
		}
		// when function ends, closes cursor
		defer cursor.Close(ctx)

		// defines genres
		var genres []models.Genre
		if err := cursor.All(ctx, &genres); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// uses json to write genres and ok status
		c.JSON(http.StatusOK, genres)

	}
}

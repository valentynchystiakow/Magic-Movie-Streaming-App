// marks file as part of database package
package database

// imports packages
import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// creates function that connects to database, marks it as mongo.Client because it works with mongo database and returns mongo.Client
func Connect() *mongo.Client {

	// loads environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: unable to find .env file")
	}

	// gets mongodb uri from env
	MongoDb := os.Getenv("MONGODB_URI")

	// if mongodb uri is empty, throws error
	if MongoDb == "" {
		log.Fatal("MONGODB_URI is not set")
	}
	fmt.Println("MongoDB URI: ", MongoDb)

	// sets client options
	clientOptions := options.Client().ApplyURI(MongoDb)

	// connects to database using client options
	client, err := mongo.Connect(clientOptions)
	// shows error if occurs
	if err != nil {
		return nil
	}

	return client
}

// creates function that opens collection from database
func OpenCollection(collectionName string, client *mongo.Client) *mongo.Collection {
	//  loads environment variables
	err := godotenv.Load(".env")
	// throws error if occurs
	if err != nil {
		log.Println("Warning: unable to find .env file")
	}
	// gets database name
	databaseName := os.Getenv("DATABASE_NAME")
	fmt.Println("Database name is: ", databaseName)

	// loads relevant collection from database
	collection := client.Database(databaseName).Collection(collectionName)

	// if collection is empty, returns nil
	if collection == nil {
		return nil
	}

	return collection
}

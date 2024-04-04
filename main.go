package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"net/http"
	"github.com/gorilla/mux"
	"swift-buddies-backend/models"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreatePostEndPoint(w http.ResponseWriter, r *http.Request) { 
	coll := client.Database("sample_mflix").Collection("posts")

	defer r.Body.Close()

	var newPost post.PostInput
	if err := json.NewDecoder(r.Body).Decode(&newPost); err != nil {
		// respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	result, err := coll.InsertOne(context.TODO(), newPost)
	if err != nil {
		panic(err)
	}

	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	
	fmt.Printf("%s\n", jsonData)
}

func GetPostEndPoint(w http.ResponseWriter, r *http.Request) { 
	coll := client.Database("sample_mflix").Collection("posts")

    var result []post.PostModel // Assuming you have imported the package containing the Post struct

    // Find all documents in the collection
    cursor, err := coll.Find(context.Background(), bson.D{})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer cursor.Close(context.Background())

    // Decode each document into the result slice
    for cursor.Next(context.Background()) {
        var post post.PostModel
        if err := cursor.Decode(&post); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        result = append(result, post)
    }

    // Check for errors encountered during cursor iteration
    if err := cursor.Err(); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Encode the result slice to JSON and write it as the response
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(result); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

var client *mongo.Client

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	mongodb, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	client = mongodb

	r := mux.NewRouter()

	r.HandleFunc("/posts", GetPostEndPoint).Methods("GET")
	r.HandleFunc("/create/post", CreatePostEndPoint).Methods("POST")

	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}

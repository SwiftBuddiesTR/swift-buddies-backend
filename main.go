package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	post "swift-buddies-backend/models"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreatePostEndPoint(w http.ResponseWriter, r *http.Request) {
	// Check for request body
	if r.Body == nil || r.ContentLength == 0 {
		http.Error(w, "Request must have a body", http.StatusBadRequest)
		return
	}

	coll := client.Database("sample_mflix").Collection("posts")

	defer r.Body.Close()

	var newPost post.PostInput // Adjusted for corrected struct name
	if err := json.NewDecoder(r.Body).Decode(&newPost); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	result, err := coll.InsertOne(context.TODO(), newPost)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert post: %v", err), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal response: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData) // Properly respond with the inserted document information
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

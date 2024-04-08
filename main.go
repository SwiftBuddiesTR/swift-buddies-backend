package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	models "swift-buddies-backend/models"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreatePostEndPoint(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil || r.ContentLength == 0 {
		http.Error(w, "Request must have a body", http.StatusBadRequest)
		return
	}

	coll := client.Database("sample_mflix").Collection("posts")
	defer r.Body.Close()

	var newPost models.PostInput
	if err := json.NewDecoder(r.Body).Decode(&newPost); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Insert post into the database
	result, err := coll.InsertOne(context.TODO(), newPost)
	if err != nil {
		http.Error(w, "Failed to insert post", http.StatusInternalServerError)
		log.Printf("Failed to insert post: %v", err)
		return
	}

	// Fetch the inserted post using the InsertedID
	var savedPost models.PostModel
	if err := coll.FindOne(context.TODO(), bson.M{"_id": result.InsertedID}).Decode(&savedPost); err != nil {
		http.Error(w, "Failed to fetch saved post", http.StatusInternalServerError)
		log.Printf("Failed to fetch saved post: %v", err)
		return
	}

	// Respond with the saved post
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(savedPost); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Failed to encode saved post: %v", err)
		return
	}
}

func GetPostEndPoint(w http.ResponseWriter, r *http.Request) {
	coll := client.Database("sample_mflix").Collection("posts")

	var result []models.PostModel // Assuming you have imported the package containing the Post struct

	// Find all documents in the collection
	cursor, err := coll.Find(context.Background(), bson.D{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	// Decode each document into the result slice
	for cursor.Next(context.Background()) {
		var post models.PostModel
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

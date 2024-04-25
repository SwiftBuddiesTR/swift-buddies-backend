package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "google.golang.org/api/idtoken"
)

// User struct to represent user information
type User struct {
	RegisterType string `json:"registerType"`
	AccessToken  string `json:"accessToken"`
}

type SimpleResponse struct {
	Message string `json:"message"`
}

// App struct to encapsulate application state
type App struct {
	Collection *mongo.Collection
}

// NewApp creates a new instance of the application
func NewApp() (*App, error) {
	// Read MongoDB connection URI from environment variable
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		return nil, errors.New("MONGODB_URI environment variable not set")
	}

	// Connect to MongoDB
	ctx := context.Background()
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Access database and collection
	collection := client.Database("buddiesapp").Collection("users")

	return &App{
		Collection: collection,
	}, nil
}

// Handler for registering a new user
func (a *App) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse JSON request body into User struct
	decoder := json.NewDecoder(r.Body)
	var user User
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if user.RegisterType != "google" && user.RegisterType != "apple" {
		http.Error(w, "Invalid registerType", http.StatusBadRequest)
		return
	}

	if user.AccessToken == "" {
		http.Error(w, "accessToken is required", http.StatusBadRequest)
		return
	}

	if user.RegisterType == "google" {
		bearerToken := user.AccessToken

		req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
		req.Header.Set("Authorization", "Bearer "+bearerToken)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusUnauthorized {

			response := SimpleResponse{
				Message: "Unauthorized",
			}

			jsonResponse, err := json.Marshal(response)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonResponse)
			return
		}

		// Read and print response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}
		fmt.Println(string(body))

		// TODO:
		// Check if the user is already registered
		// After checking, if the user is not registered, insert the user into MongoDB collection
		// If the user is already registered, return a token and refreshToken
	}

	// Insert user into MongoDB collection
	// _, err := a.Collection.InsertOne(context.Background(), user)
	// if err != nil {
	//     http.Error(w, "Failed to register user", http.StatusInternalServerError)
	//     return
	// }

	// Return a success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	// Create a new instance of the application
	app, err := NewApp()
	if err != nil {
		log.Fatal(err)
	}

	// Define HTTP routes
	http.HandleFunc("/register", app.RegisterUserHandler)

	// Start the HTTP server
	fmt.Println("Server running on localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

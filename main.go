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
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewApp() (*App, error) {
	godotenv.Load()

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

func GenerateToken(claims jwt.MapClaims, expiration time.Time, secretKey []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Handler for registering a new user
func (a *App) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var user RegisterUser
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
		fmt.Println("body", string(body))

		var userCredentials GoogleUserCredentials
		err = json.Unmarshal(body, &userCredentials)
		if err != nil {
			fmt.Println("Error unmarshalling response:", err)
			return
		}

		secretKey := os.Getenv("SECRET_KEY")
		if secretKey == "" {
			fmt.Println("SECRET_KEY environment variable not set")
			return
		}

		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["sub"] = "user"
		claims["exp"] = time.Now().Add(1 * time.Hour).Unix()

		// Convert secretKey to []byte
		secretKeyBytes := []byte(secretKey)

		// Generate token
		tokenString, err := token.SignedString(secretKeyBytes)
		if err != nil {
			fmt.Println("Error generating token:", err)
			return
		}

		// response tokenString with json format
		response := SimpleResponse{
			Message: tokenString,
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)

		// TODO:
		// Check if the user is already registered

		filter := bson.M{"email": "dogukaankilicarslan@gmail.com"}

		// Try to find the user
		var existingUser User
		fmt.Println("finding")
		err = a.Collection.FindOne(context.Background(), filter).Decode(&existingUser)
		fmt.Println("okey")
		if err == mongo.ErrNoDocuments {
			// If user not found, insert the new user
			uuid := uuid.New()

			newUser := User{
				ID:           uuid.String(), 
				Name:         string(userCredentials.Name),
				Email:        string(userCredentials.Email),
				Picture:      string(userCredentials.Picture),
				Token:        tokenString,
				RefreshToken: "TODO",
			}

			// Insert the new user into the collection
			_, err := a.Collection.InsertOne(context.Background(), newUser)
			if err != nil {
				log.Fatal(err)
			}

			// Return the token and refresh token
			fmt.Println("Token:", newUser.Token)
			fmt.Println("Refresh Token:", newUser.RefreshToken)
			var registerReturn RegisterReturn
			registerReturn.Token = newUser.Token
			registerReturn.RefreshToken = newUser.RefreshToken
			// response registerReturn with json format
			jsonResponse, err := json.Marshal(registerReturn)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonResponse)
			return			
		} else if err != nil {
			log.Fatal(err)
		} else {
			// If user found, return the existing token and refresh token
			fmt.Println("Token:", existingUser.Token)
			fmt.Println("Refresh Token:", existingUser.RefreshToken)
			var registerReturn RegisterReturn
			registerReturn.Token = existingUser.Token
			registerReturn.RefreshToken = existingUser.RefreshToken
			// response registerReturn with json format
			jsonResponse, err := json.Marshal(registerReturn)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonResponse)
			return		
		}

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

package main

import "go.mongodb.org/mongo-driver/mongo"

type GoogleUserCredentials struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

type User struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Picture      string `json:"picture"`
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type RegisterReturn struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type RegisterUser struct {
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

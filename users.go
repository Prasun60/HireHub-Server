package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"bytes"
	"strconv"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
	"os"
	
	
)

type user struct {
	FullName     string `json:"fullname"`
	UserName     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	IsRecruiter  bool   `json:"is_recruiter"`
	Resume       string `json:"resume"`
	ProfileImage string `json:"profile_image"`
	Address      string `json:"address"`
	Bio          string `json:"bio"`
	Description  string `json:"description"`
	Position     string `json:"position"`
	Skills       string `json:"skills"`
	IsPremium    bool   `json:"is_premium"`
	Portfolio    string `json:"portfolio"`
	
}



func Register_Recruiter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var person user
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Ensure IsRecruiter is set to true
	person.IsRecruiter = true
	person.IsPremium = false

	fmt.Printf("User to be inserted: %+v\n", person) // Debug print

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(person.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	person.Password = string(hashedPassword)

	// Generate a unique username
	person.UserName = generateUniqueUsername(person.FullName, "-recruiter")

	// Insert user into the database
	insertUserQuery := `
		INSERT INTO users (fullname, username, email, password, is_recruiter, resume, profile_image, address, bio, description, position, portfolio, is_premium, skills)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id
	`

	var userID int
	err = db.QueryRow(insertUserQuery, person.FullName, person.UserName, person.Email, person.Password, person.IsRecruiter, person.Resume, person.ProfileImage, person.Address, person.Bio, person.Description, person.Position, person.Portfolio, person.IsPremium, person.Skills).Scan(&userID)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Failed to insert recruiter", http.StatusInternalServerError)
		return
	}

	fmt.Println("Inserted a single document with ID: ", userID)
	json.NewEncoder(w).Encode(map[string]int{"user_id": userID})
}

func Register_Normal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var person user
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Ensure IsRecruiter is set to true
	person.IsRecruiter = false
	person.IsPremium = false

	fmt.Printf("User to be inserted: %+v\n", person) // Debug print

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(person.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	person.Password = string(hashedPassword)

	// Generate a unique username
	person.UserName = generateUniqueUsername(person.FullName, "-user")

	// Insert user into the database
	insertUserQuery := `
		INSERT INTO users (fullname, username, email, password, is_recruiter, resume, profile_image, address, bio, description, position, portfolio, is_premium, skills)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id
	`

	var userID int
	err = db.QueryRow(insertUserQuery, person.FullName, person.UserName, person.Email, person.Password, person.IsRecruiter, person.Resume, person.ProfileImage, person.Address, person.Bio, person.Description, person.Position, person.Portfolio, person.IsPremium, person.Skills).Scan(&userID)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Failed to insert user", http.StatusInternalServerError)
		return
	}

	fmt.Println("Inserted a single document with ID: ", userID)
	json.NewEncoder(w).Encode(map[string]int{"user_id": userID})
}

func generateUniqueUsername(fullName string, suffix string) string {
	if len(fullName) > 4 {
		fullName = fullName[:4]
	}
	return fmt.Sprintf("%s_%d%s", fullName, time.Now().Unix(), suffix)
}


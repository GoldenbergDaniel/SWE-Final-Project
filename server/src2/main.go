package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

// Initialize the database connection and create the users table with additional fields
func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Drop the table if it exists and recreate it
	_, err = db.Exec(`
        DROP TABLE IF EXISTS users;
        CREATE TABLE users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            first_name TEXT NOT NULL,
            last_name TEXT NOT NULL,
            email TEXT UNIQUE NOT NULL,
            username TEXT UNIQUE NOT NULL,
            password TEXT NOT NULL
        )
    `)
	if err != nil {
		log.Fatal("Error creating users table:", err)
	}
	fmt.Println("Connected to database and ensured users table exists.")
}

func main() {
	// Set up logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	initDB()
	defer db.Close()

	r := mux.NewRouter()

	// Add OPTIONS handler for CORS preflight requests
	r.HandleFunc("/{path}", PreflightHandler).Methods("OPTIONS")
	r.HandleFunc("/signup", PostSignup).Methods("POST")
	r.HandleFunc("/login", PostLogin).Methods("POST")

	// Add middleware to log all requests
	r.Use(loggingMiddleware)

	fmt.Println("Server is running on localhost:5174")
	log.Fatal(http.ListenAndServe(":5174", r))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received %s request to %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func PreflightHandler(w http.ResponseWriter, r *http.Request) {
	SetCors(w.Header())
	w.WriteHeader(http.StatusOK)
}

func PostSignup(w http.ResponseWriter, r *http.Request) {
	SetCors(w.Header())
	log.Printf("Content-Type: %s", r.Header.Get("Content-Type"))

	var credentials struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Username  string `json:"username"`
		Password  string `json:"password"`
	}

	// Read the body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	log.Printf("Received body: %s", string(body))

	// Decode the JSON
	if err := json.Unmarshal(body, &credentials); err != nil {
		log.Printf("JSON decode error: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate input
	if credentials.Username == "" || credentials.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Hash the password for secure storage
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Password hashing error: %v", err)
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Insert the user into the database
	result, err := db.Exec(`
        INSERT INTO users (first_name, last_name, email, username, password)
        VALUES (?, ?, ?, ?, ?)
    `, credentials.FirstName, credentials.LastName, credentials.Email, credentials.Username, string(hashedPassword))

	if err != nil {
		log.Printf("Database insert error: %v", err)
		http.Error(w, "Failed to insert user", http.StatusInternalServerError)
		return
	}

	// Get the ID of the inserted user
	id, _ := result.LastInsertId()

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User registered successfully",
		"user_id": id,
	})
}

func PostLogin(w http.ResponseWriter, r *http.Request) {
	SetCors(w.Header())

	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		log.Printf("JSON decode error: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Retrieve the hashed password from the database
	var hashedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", credentials.Username).Scan(&hashedPassword)
	if err != nil {
		log.Printf("Database query error: %v", err)
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Compare the provided password with the stored hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(credentials.Password)); err != nil {
		log.Printf("Password comparison error: %v", err)
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful",
	})
}

func SetCors(header http.Header) {
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, OPTIONS")
	header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

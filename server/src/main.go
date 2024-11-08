package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"github.com/rs/cors"
)

var db *sql.DB

// Initialize the database connection and create the users table if it doesn't exist
func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Create the users table if it doesn't exist
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
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

	// Add your routes
	r.HandleFunc("/signup", PostSignup).Methods("POST")
	r.HandleFunc("/login", PostLogin).Methods("POST")
	r.HandleFunc("/logout", Logout).Methods("POST")
	r.HandleFunc("/protected", AuthMiddleware(ProtectedHandler)).Methods("GET")

	// Create a new CORS handler
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Wrap your router with the CORS handler
	handler := c.Handler(r)

	// Use the new handler instead of r
	log.Fatal(http.ListenAndServe(":5174", handler))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received %s request to %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func isValidEmail(email string) bool {
	// Basic regex for validating an email address
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
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
	if credentials.FirstName == "" || credentials.LastName == "" {
		http.Error(w, "First Name and Last Name are required", http.StatusBadRequest)
		return
	}

	if credentials.Email == "" || !isValidEmail(credentials.Email) {
		http.Error(w, "A valid email is required", http.StatusBadRequest)
		return
	}

	if credentials.Username == "" || credentials.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Check for existing email
	var exists int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", credentials.Email).Scan(&exists)
	if err != nil {
		log.Printf("Database query error: %v", err)
		http.Error(w, "Error checking email", http.StatusInternalServerError)
		return
	}
	if exists > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Email already exists",
		})
		return
	}

	// Check for existing username
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", credentials.Username).Scan(&exists)
	if err != nil {
		log.Printf("Database query error: %v", err)
		http.Error(w, "Error checking username", http.StatusInternalServerError)
		return
	}
	if exists > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Username already exists",
		})
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

	// Create a session cookie
	sessionToken := generateSessionToken()
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful",
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	SetCors(w.Header())

	// Clear the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logged out successfully",
	})
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		SetCors(w.Header())

		cookie, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// Here you would typically validate the session token
		// For this example, we'll just check if it exists
		if cookie.Value == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "This is a protected route",
	})
}

func generateSessionToken() string {
	// For this example, we'll just use a timestamp
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func SetCors(header http.Header) {
	header.Set("Access-Control-Allow-Origin", "http://localhost:5173")
	header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	header.Set("Access-Control-Allow-Credentials", "true")
}

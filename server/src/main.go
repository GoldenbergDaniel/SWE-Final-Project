package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

// Initialize the database connection and create the users table if it doesn't exist
func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./data.db")
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
	}

	if err = db.Ping(); err != nil {
		fmt.Println("Failed to ping database:", err)
	}

	// Create the users table if it doesn't exist
	_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS users (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					first_name TEXT NOT NULL,
					last_name TEXT NOT NULL,
					email TEXT UNIQUE NOT NULL,
					username TEXT UNIQUE NOT NULL,
					password TEXT NOT NULL,
					balance INTEGER NOT NULL
			)
	`)
	if err != nil {
		fmt.Println("Error creating users table:", err)
	}
	fmt.Println("Connected to database and ensured users table exists.")
}

func main() {
	initDB()
	defer db.Close()

	r := mux.NewRouter()

	// Add your routes
	r.HandleFunc("/signup", PostSignup).Methods("POST")
	r.HandleFunc("/login", PostLogin).Methods("POST")
	r.HandleFunc("/logout", Logout).Methods("POST")
	r.HandleFunc("/protected", AuthMiddleware(ProtectedHandler)).Methods("GET")
	r.HandleFunc("/userdata", GetUserData).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)
	fmt.Println(http.ListenAndServe(":5174", handler))
}

func IsEmailValid(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func PreflightHandler(w http.ResponseWriter, r *http.Request) {
	// SetCors(w.Header())
	w.WriteHeader(http.StatusOK)
}

func GetUserData(w http.ResponseWriter, r *http.Request) {
	// SetCors(w.Header())

	cookie, err := r.Cookie("session_token")
	if err != nil {
		fmt.Println("Error! No cookie provided.")
		return
	}

	var balance int
	err = db.QueryRow("SELECT balance FROM users WHERE username = ?", cookie.Value).Scan(&balance)
	if err != nil {
		fmt.Println("User not found!")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"balance": balance,
	})
}

func PostSignup(w http.ResponseWriter, r *http.Request) {
	// SetCors(w.Header())
	fmt.Printf("Content-Type: %s", r.Header.Get("Content-Type"))

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
		fmt.Printf("Error reading body: %v", err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received body: %s", string(body))

	// Decode the JSON
	if err := json.Unmarshal(body, &credentials); err != nil {
		fmt.Printf("JSON decode error: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate input
	if credentials.FirstName == "" || credentials.LastName == "" {
		http.Error(w, "First Name and Last Name are required", http.StatusBadRequest)
		return
	}

	if credentials.Email == "" || !IsEmailValid(credentials.Email) {
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
		fmt.Printf("Database query error: %v", err)
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
		fmt.Printf("Database query error: %v", err)
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
		fmt.Printf("Password hashing error: %v", err)
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Insert the user into the database
	result, err := db.Exec(`
        INSERT INTO users (first_name, last_name, email, username, password, balance)
        VALUES (?, ?, ?, ?, ?, ?)
    `, credentials.FirstName, credentials.LastName, credentials.Email, credentials.Username, string(hashedPassword), 10000)

	if err != nil {
		fmt.Printf("Database insert error: %v", err)
		http.Error(w, "Failed to insert user", http.StatusInternalServerError)
		return
	}

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
	// SetCors(w.Header())

	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		fmt.Printf("JSON decode error: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Retrieve the hashed password from the database
	var hashedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", credentials.Username).Scan(&hashedPassword)
	if err != nil {
		fmt.Printf("Database query error: %v", err)
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Compare the provided password with the stored hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(credentials.Password)); err != nil {
		fmt.Printf("Password comparison error: %v", err)
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	// Create a session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    credentials.Username,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false,
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
	// SetCors(w.Header())

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
		// SetCors(w.Header())

		cookie, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// Check if cookie exists
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

// func SetCors(header http.Header) {
// 	header.Set("Access-Control-Allow-Origin", "http://localhost:5173")
// 	header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// 	header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
// 	header.Set("Access-Control-Allow-Credentials", "true")
// }

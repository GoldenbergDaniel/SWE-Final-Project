package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/robfig/cron/v3"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var stockCache map[string]StockPrice

type StockPrice struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
	Time   string  `json:"time"`
}

func getCookieValue(r *http.Request, cookieName string) string {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	tables := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			balance REAL NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS trades (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			symbol TEXT NOT NULL,
			quantity INTEGER NOT NULL,
			price REAL NOT NULL,
			trade_type TEXT NOT NULL,
			trade_date DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS portfolio (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			symbol TEXT NOT NULL,
			quantity INTEGER NOT NULL,
			average_price REAL NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id),
			UNIQUE(user_id, symbol)
		)`,
		`CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			symbol TEXT NOT NULL,
			quantity INTEGER NOT NULL,
			trade_type TEXT NOT NULL,
			rationale TEXT,
			trade_date DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS posts_likes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			post_id INTEGER NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (post_id) REFERENCES posts(id),
			UNIQUE(user_id, post_id)
		)`,
		`CREATE TABLE IF NOT EXISTS daily_stock_prices (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			symbol TEXT NOT NULL,
			price REAL NOT NULL,
			updated_at DATETIME NOT NULL,
			UNIQUE(symbol, updated_at)
		)`,
	}

	for _, table := range tables {
		_, err := db.Exec(table)
		if err != nil {
			log.Printf("Error creating table: %v", err)
		}
	}

	log.Println("Connected to database and ensured all tables exist.")
}

func main() {
	initDB()
	defer db.Close()

	stockCache = make(map[string]StockPrice)

	r := mux.NewRouter()

	r.HandleFunc("/signup", PostSignup).Methods("POST")
	r.HandleFunc("/login", PostLogin).Methods("POST")
	r.HandleFunc("/logout", Logout).Methods("POST")
	r.HandleFunc("/protected", AuthMiddleware(ProtectedHandler)).Methods("GET")
	r.HandleFunc("/userdata", GetUserData).Methods("GET")
	r.HandleFunc("/stock-price", GetStockPrice).Methods("GET")
	r.HandleFunc("/trade", AuthMiddleware(MakeTrade)).Methods("POST")
	r.HandleFunc("/portfolio-value", AuthMiddleware(GetPortfolioValue)).Methods("GET")
	r.HandleFunc("/historical-prices", AuthMiddleware(GetHistoricalPrices)).Methods("GET")
	r.HandleFunc("/leaderboard", AuthMiddleware(GetLeaderboard)).Methods("GET")
	r.HandleFunc("/posts", AuthMiddleware(GetPosts)).Methods("GET")
	r.HandleFunc("/like/{id}", AuthMiddleware(ToggleLike)).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	startStockPriceUpdateJob()

	fmt.Println(http.ListenAndServe(":5174", handler))
}

func startStockPriceUpdateJob() {
	c := cron.New()
	c.AddFunc("0 9 * * *", updateDailyStockPrices) // Run every day at 9:00 AM
	c.Start()
}

func updateDailyStockPrices() {
	fmt.Printf("Updating daily stock prices at %s\n", time.Now().Format(time.RFC3339))

	symbols, err := getUniqueSymbolsInPortfolios()
	if err != nil {
		fmt.Println("Error getting unique symbols:", err)
		return
	}

	for _, symbol := range symbols {
		price, err := fetchStockPrice(symbol)
		if err != nil {
			fmt.Printf("Error fetching price for %s: %v\n", symbol, err)
			continue
		}

		_, err = db.Exec(`
			INSERT INTO daily_stock_prices (symbol, price, updated_at)
			VALUES (?, ?, datetime('now'))
		`, symbol, price)
		if err != nil {
			fmt.Printf("Error storing daily price for %s: %v\n", symbol, err)
		} else {
			fmt.Printf("Updated daily price for %s: $%.2f\n", symbol, price)
		}
	}

	fmt.Println("Daily stock price update completed")
}

func getUniqueSymbolsInPortfolios() ([]string, error) {
	rows, err := db.Query("SELECT DISTINCT symbol FROM portfolio")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var symbols []string
	for rows.Next() {
		var symbol string
		if err := rows.Scan(&symbol); err != nil {
			return nil, err
		}
		symbols = append(symbols, symbol)
	}

	return symbols, nil
}

func fetchStockPrice(symbol string) (float64, error) {
	apiKey := "J585VGMES541XQW2"
	url := fmt.Sprintf("https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=%s&apikey=%s", symbol, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	globalQuote, ok := result["Global Quote"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("invalid response from API")
	}

	price, err := strconv.ParseFloat(globalQuote["05. price"].(string), 64)
	if err != nil {
		return 0, err
	}

	return price, nil
}

func IsEmailValid(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func GetUserData(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		fmt.Println("Error! No cookie provided.")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var balance float64
	err = db.QueryRow("SELECT balance FROM users WHERE username = ?", cookie.Value).Scan(&balance)
	if err != nil {
		fmt.Println("Username or Password Incorrect!")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{
		"balance": balance,
	})
}

func PostSignup(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Username  string `json:"username"`
		Password  string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

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

	var exists int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", credentials.Email).Scan(&exists)
	if err != nil {
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

	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", credentials.Username).Scan(&exists)
	if err != nil {
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	result, err := db.Exec(`
		INSERT INTO users (first_name, last_name, email, username, password, balance)
		VALUES (?, ?, ?, ?, ?, ?)
	`, credentials.FirstName, credentials.LastName, credentials.Email, credentials.Username, string(hashedPassword), 10000.0)

	if err != nil {
		http.Error(w, "Failed to insert user", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User registered successfully",
		"user_id": id,
	})
}

func PostLogin(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var hashedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", credentials.Username).Scan(&hashedPassword)
	if err != nil {
		http.Error(w, "Username or Password Incorrect", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(credentials.Password)); err != nil {
		http.Error(w, "Username or Password Incorrect", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    credentials.Username,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful",
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
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
		cookie, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

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

func GetStockPrice(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		http.Error(w, "Symbol is required", http.StatusBadRequest)
		return
	}

	if cachedPrice, ok := stockCache[symbol]; ok {
		parsedTime, err := time.Parse(time.RFC3339, cachedPrice.Time)
		if err == nil && time.Since(parsedTime) < 5*time.Minute {
			json.NewEncoder(w).Encode(cachedPrice)
			return
		}
	}

	price, err := fetchStockPrice(symbol)
	if err != nil {
		http.Error(w, "Failed to fetch stock price", http.StatusInternalServerError)
		return
	}

	stockPrice := StockPrice{
		Symbol: symbol,
		Price:  price,
		Time:   time.Now().Format(time.RFC3339)}

	stockCache[symbol] = stockPrice

	json.NewEncoder(w).Encode(stockPrice)
}

func getUserIdFromSession(r *http.Request) int {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return 0
	}

	var userId int
	err = db.QueryRow("SELECT id FROM users WHERE username = ?", cookie.Value).Scan(&userId)
	if err != nil {
		return 0
	}

	return userId
}

func MakeTrade(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var tradeReq struct {
		Symbol    string `json:"symbol"`
		Quantity  int    `json:"quantity"`
		TradeType string `json:"trade_type"`
		Rationale string `json:"rationale"`
	}

	if err := json.NewDecoder(r.Body).Decode(&tradeReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if tradeReq.Quantity <= 0 {
		http.Error(w, "Quantity must be greater than 0", http.StatusBadRequest)
		return
	}

	cookie, _ := r.Cookie("session_token")
	username := cookie.Value

	stockPrice, err := fetchStockPrice(tradeReq.Symbol)
	if err != nil {
		http.Error(w, "Failed to get stock price", http.StatusInternalServerError)
		return
	}

	totalCost := float64(tradeReq.Quantity) * stockPrice

	var balance float64
	var userId int
	err = db.QueryRow("SELECT id, balance FROM users WHERE username = ?", username).Scan(&userId, &balance)
	if err != nil {
		http.Error(w, "Failed to get user data", http.StatusInternalServerError)
		return
	}

	if tradeReq.TradeType == "buy" && balance < totalCost {
		http.Error(w, "Insufficient balance", http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	var newBalance float64
	if tradeReq.TradeType == "buy" {
		newBalance = balance - totalCost
	} else {
		newBalance = balance + totalCost
	}

	_, err = tx.Exec("UPDATE users SET balance = ? WHERE id = ?", newBalance, userId)
	if err != nil {
		http.Error(w, "Failed to update balance", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`
		INSERT INTO trades (user_id, symbol, quantity, price, trade_type)
		VALUES (?, ?, ?, ?, ?)
	`, userId, tradeReq.Symbol, tradeReq.Quantity, stockPrice, tradeReq.TradeType)
	if err != nil {
		http.Error(w, "Failed to record trade", http.StatusInternalServerError)
		return
	}

	var currentQuantity int
	err = tx.QueryRow("SELECT quantity FROM portfolio WHERE user_id = ? AND symbol = ?", userId, tradeReq.Symbol).Scan(&currentQuantity)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Failed to get current portfolio quantity", http.StatusInternalServerError)
		return
	}

	var newQuantity int
	if tradeReq.TradeType == "buy" {
		newQuantity = currentQuantity + tradeReq.Quantity
	} else {
		newQuantity = currentQuantity - tradeReq.Quantity
	}

	if newQuantity < 0 {
		http.Error(w, "Insufficient shares to sell", http.StatusBadRequest)
		return
	}

	if newQuantity == 0 {
		_, err = tx.Exec("DELETE FROM portfolio WHERE user_id = ? AND symbol = ?", userId, tradeReq.Symbol)
	} else {
		_, err = tx.Exec(`
			INSERT OR REPLACE INTO portfolio (user_id, symbol, quantity, average_price)
			VALUES (?, ?, ?, (SELECT COALESCE(
				(SELECT (average_price * quantity + ? * ?) / (quantity + ?)
				FROM portfolio WHERE user_id = ? AND symbol = ?),
				?
			)))
		`, userId, tradeReq.Symbol, newQuantity, tradeReq.Quantity, stockPrice, tradeReq.Quantity, userId, tradeReq.Symbol, stockPrice)
	}

	if err != nil {
		http.Error(w, "Failed to update portfolio", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`
	INSERT INTO posts (user_id, symbol, quantity, trade_type, rationale, trade_date)
	VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
`, userId, tradeReq.Symbol, tradeReq.Quantity, tradeReq.TradeType, tradeReq.Rationale)

	if err != nil {
		http.Error(w, "Failed to create post for trade", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":     "Trade successful",
		"new_balance": newBalance,
	})
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	userId := getUserIdFromSession(r)

	rows, err := db.Query(`
		SELECT p.id, u.username, p.symbol, p.quantity, p.trade_type, p.rationale, p.trade_date,
			   (SELECT COUNT(*) FROM posts_likes WHERE post_id = p.id) AS likes_count,
			   CASE WHEN pl.user_id IS NOT NULL THEN 1 ELSE 0 END AS liked_by_user
		FROM posts p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN posts_likes pl ON p.id = pl.post_id AND pl.user_id = ?
		ORDER BY p.trade_date DESC
		LIMIT 50
	`, userId)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []map[string]interface{}
	for rows.Next() {
		var postId int
		var username, symbol, tradeType, rationale string
		var quantity, likesCount int
		var likedByUser bool
		var tradeDate time.Time
		err := rows.Scan(&postId, &username, &symbol, &quantity, &tradeType, &rationale, &tradeDate, &likesCount, &likedByUser)
		if err != nil {
			http.Error(w, "Failed to scan post row", http.StatusInternalServerError)
			return
		}

		posts = append(posts, map[string]interface{}{
			"id":            postId,
			"username":      username,
			"symbol":        symbol,
			"quantity":      quantity,
			"trade_type":    tradeType,
			"rationale":     rationale,
			"trade_date":    tradeDate,
			"likes":         likesCount,
			"liked_by_user": likedByUser,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func ToggleLike(w http.ResponseWriter, r *http.Request) {
	postId := mux.Vars(r)["id"]
	userId := getUserIdFromSession(r)

	var liked bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM posts_likes WHERE user_id = ? AND post_id = ?)", userId, postId).Scan(&liked)
	if err != nil {
		http.Error(w, "Failed to check like status", http.StatusInternalServerError)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	if liked {
		_, err = tx.Exec("DELETE FROM posts_likes WHERE user_id = ? AND post_id = ?", userId, postId)
	} else {
		_, err = tx.Exec("INSERT INTO posts_likes (user_id, post_id) VALUES (?, ?)", userId, postId)
	}

	if err != nil {
		http.Error(w, "Failed to toggle like", http.StatusInternalServerError)
		return
	}

	var likesCount int
	err = tx.QueryRow("SELECT COUNT(*) FROM posts_likes WHERE post_id = ?", postId).Scan(&likesCount)
	if err != nil {
		http.Error(w, "Failed to get updated like count", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"likes":         likesCount,
		"liked_by_user": !liked,
	})
}

func GetPortfolioValue(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		fmt.Println("Error getting session cookie:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	username := cookie.Value

	var userId int
	var email string
	var balance float64
	err = db.QueryRow("SELECT id, email, balance FROM users WHERE username = ?", username).Scan(&userId, &email, &balance)
	if err != nil {
		fmt.Println("Error querying user data:", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	rows, err := db.Query(`
		SELECT p.symbol, p.quantity, p.average_price, COALESCE(dsp.price, p.average_price) as current_price
		FROM portfolio p
		LEFT JOIN (
			SELECT symbol, price
			FROM daily_stock_prices
			WHERE updated_at = (SELECT MAX(updated_at) FROM daily_stock_prices)
		) dsp ON p.symbol = dsp.symbol
		WHERE p.user_id = ?
	`, userId)
	if err != nil {
		fmt.Println("Error querying portfolio data:", err)
		http.Error(w, "Failed to fetch portfolio", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var totalValue float64
	portfolio := make(map[string]map[string]interface{})

	for rows.Next() {
		var symbol string
		var quantity int
		var averagePrice, currentPrice float64
		err := rows.Scan(&symbol, &quantity, &averagePrice, &currentPrice)
		if err != nil {
			fmt.Println("Error scanning portfolio row:", err)
			http.Error(w, "Failed to scan portfolio row", http.StatusInternalServerError)
			return
		}

		marketValue := float64(quantity) * currentPrice
		totalValue += marketValue

		portfolio[symbol] = map[string]interface{}{
			"quantity":     quantity,
			"averagePrice": averagePrice,
			"currentPrice": currentPrice,
			"marketValue":  marketValue,
			"profitLoss":   marketValue - (float64(quantity) * averagePrice),
		}
	}

	totalValue += balance

	response := map[string]interface{}{
		"username":   username,
		"email":      email,
		"balance":    balance,
		"totalValue": totalValue,
		"portfolio":  portfolio,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Println("Error encoding response:", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func GetHistoricalPrices(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		http.Error(w, "Symbol is required", http.StatusBadRequest)
		return
	}

	days := 30 // Default to 30 days
	if daysParam := r.URL.Query().Get("days"); daysParam != "" {
		if parsedDays, err := strconv.Atoi(daysParam); err == nil {
			days = parsedDays
		}
	}

	rows, err := db.Query(`
		SELECT date, price
		FROM historical_prices
		WHERE symbol = ?
		AND date >= date('now', '-' || ? || ' days')
		ORDER BY date ASC
	`, symbol, days)
	if err != nil {
		http.Error(w, "Failed to fetch historical prices", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var prices []map[string]interface{}
	for rows.Next() {
		var date string
		var price float64
		if err := rows.Scan(&date, &price); err != nil {
			http.Error(w, "Failed to scan historical price row", http.StatusInternalServerError)
			return
		}
		prices = append(prices, map[string]interface{}{
			"date":  date,
			"price": price,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"symbol": symbol,
		"prices": prices,
	})
}

func GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
        SELECT u.id, u.username, u.balance
        FROM users u
        ORDER BY u.balance DESC
        LIMIT 10
    `)
	if err != nil {
		http.Error(w, "Failed to fetch leaderboard data", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var leaderboard []map[string]interface{}
	for rows.Next() {
		var userId int
		var username string
		var balance float64
		err := rows.Scan(&userId, &username, &balance)
		if err != nil {
			http.Error(w, "Failed to scan leaderboard row", http.StatusInternalServerError)
			return
		}

		portfolioValue, err := getPortfolioValue(userId)
		if err != nil {
			http.Error(w, "Failed to get portfolio value", http.StatusInternalServerError)
			return
		}

		totalValue := balance + portfolioValue
		gainLoss := (totalValue - 10000) / 100 // Assuming initial balance was 10000

		leaderboard = append(leaderboard, map[string]interface{}{
			"username":   username,
			"totalValue": totalValue,
			"gainLoss":   gainLoss,
		})
	}

	// Sort the leaderboard by totalValue in descending order
	sort.Slice(leaderboard, func(i, j int) bool {
		return leaderboard[i]["totalValue"].(float64) > leaderboard[j]["totalValue"].(float64)
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(leaderboard)
}

// Helper function to calculate portfolio value
func getPortfolioValue(userId int) (float64, error) {
	/* OLD ONE
		rows, err := db.Query(`
	        SELECT p.quantity, COALESCE(dsp.price, p.average_price) as current_price
	        FROM portfolio p
	        LEFT JOIN (
	            SELECT symbol, price
	            FROM daily_stock_prices
	            WHERE updated_at = (SELECT MAX(updated_at) FROM daily_stock_prices)
	        ) dsp ON p.symbol = dsp.symbol
	        WHERE p.user_id = ?
	    `, userId)
	*/
	rows, err := db.Query(`
    	SELECT p.quantity, COALESCE(dsp.price, p.average_price) as current_price
    	FROM portfolio p
    	LEFT JOIN (
        	SELECT dsp.symbol, dsp.price
        	FROM daily_stock_prices dsp
        	INNER JOIN (
            	SELECT symbol, MAX(updated_at) AS latest_update
            	FROM daily_stock_prices
            	GROUP BY symbol
        	) latest_prices
        	ON dsp.symbol = latest_prices.symbol AND dsp.updated_at = latest_prices.latest_update
    	) dsp ON p.symbol = dsp.symbol
    	WHERE p.user_id = ?
	`, userId)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var totalValue float64
	for rows.Next() {
		var quantity int
		var currentPrice float64
		if err := rows.Scan(&quantity, &currentPrice); err != nil {
			return 0, err
		}
		totalValue += float64(quantity) * currentPrice
	}

	return totalValue, nil
}

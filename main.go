package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/MazenBahie/go_http_server/handlers"
	"github.com/MazenBahie/go_http_server/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	fmt.Println("Server is starting...")

	fmt.Println("Reading environment variables...")
	err := godotenv.Load(".env")

	handleFatalError(err, "Error loading .env file")
	portStr := os.Getenv("PORT")
	if portStr == "" {
		handleFatalError(errors.New("error PORT is not found in env file"), "port is not found in env file")
		return
	}
	fmt.Println("Detected port: " + portStr)

	dburl := os.Getenv("DBPath")
	if dburl == "" {
		handleFatalError(errors.New("error DBPath is not found in env file"), "DBPath is not found in env file")
		return
	}
	fmt.Println("Detected database path: " + dburl)

	fmt.Println("Connecting to database...")

	db, err := sql.Open("postgres", dburl)
	if err != nil {
		handleFatalError(err, "Error connecting to database")
	}
	defer db.Close()

	fmt.Println("Database connection established")
	fmt.Println("Creating schema if not created...")
	_, err = db.Exec(
		`
		CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    stripe_customer_id VARCHAR(255) UNIQUE,
    role VARCHAR(10) NOT NULL DEFAULT 'user'
);


CREATE TABLE IF NOT EXISTS credit_cards (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    card_number bigint NOT NULL,
    cvv VARCHAR(3) NOT NULL,
    expiration_date DATE NOT NULL,
    user_id INTEGER REFERENCES users(id)
);


CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS user_orders (
    id SERIAL PRIMARY KEY,
    order_date DATE NOT NULL,
    user_id INTEGER REFERENCES users(id),
	product_id INTEGER REFERENCES products(id)
);
		`)

	fmt.Println("Initializing dependencies...")
	// Initialize the validators
	handlers.Val = validator.New()

	fmt.Println("mapping routes...")
	r := mux.NewRouter()

	r.HandleFunc("/", handlers.HandleHome).Methods(http.MethodGet)
	r.HandleFunc("/error", handlers.HandleError)
	r.HandleFunc("/signup", handlers.HandleSignUp(db)).Methods(http.MethodPost)
	r.HandleFunc("/login", handlers.HandleLogin(db)).Methods(http.MethodPost)

	//credit card routes
	r.HandleFunc("/credit", middleware.AuthMiddleware(handlers.HandleAddCreditCard(db), false)).Methods(http.MethodPost)
	r.HandleFunc("/credit/{card_id}", middleware.AuthMiddleware(handlers.HandleDeleteCreditCard(db), false)).
		Methods(http.MethodDelete)

	//product routes
	r.HandleFunc("/products", middleware.AuthMiddleware(handlers.HandleGetProducts(db), false)).
		Methods(http.MethodGet)

	// admin products routes
	r.HandleFunc("/products/{product_id}", middleware.AuthMiddleware(handlers.HandleUpdateProduct(db), true)).
		Methods(http.MethodPut)
	r.HandleFunc("/products/{product_id}", middleware.AuthMiddleware(handlers.HandleDeleteProduct(db), true)).
		Methods(http.MethodDelete)

	r.HandleFunc("/products", middleware.AuthMiddleware(handlers.HandleCreateProduct(db), true)).
		Methods(http.MethodPost)

	r.HandleFunc("/llogin", middleware.AuthMiddleware(handlers.HandleLogin(db), true)).Methods(http.MethodPost)
	r.HandleFunc("/lllogin", middleware.AuthMiddleware(handlers.HandleLogin(db), false)).Methods(http.MethodPost)

	fmt.Println("Server is started on port => " + portStr)
	http.ListenAndServe(":"+portStr, r)
}

func handleFatalError(err error, msg string) {
	if err != nil {
		log.Fatal(msg, err)
		return
	}
}

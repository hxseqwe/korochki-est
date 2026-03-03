package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/hxseqwe/korochki-est/internal/handler"
	"github.com/hxseqwe/korochki-est/internal/repository"
	"github.com/hxseqwe/korochki-est/internal/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://korochki_user:korochki123@localhost/korochki_est?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Database connected successfully")

	sessionKey := os.Getenv("SESSION_KEY")
	if sessionKey == "" {
		sessionKey = "super-secret-key-change-in-production"
	}
	sessionStore := sessions.NewCookieStore([]byte(sessionKey))

	userRepo := repository.NewUserRepository(db)
	appRepo := repository.NewApplicationRepository(db)

	authService := service.NewAuthService(userRepo, sessionStore)
	appService := service.NewApplicationService(appRepo)

	authHandler := handler.NewAuthHandler(authService)
	appHandler := handler.NewApplicationHandler(appService, sessionStore)

	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/register", authHandler.Register).Methods("POST")
	api.HandleFunc("/login", authHandler.Login).Methods("POST")
	api.HandleFunc("/logout", authHandler.Logout).Methods("POST")

	protected := api.PathPrefix("").Subrouter()
	protected.Use(authHandler.AuthMiddleware)
	protected.HandleFunc("/applications", appHandler.GetUserApplications).Methods("GET")
	protected.HandleFunc("/applications", appHandler.CreateApplication).Methods("POST")
	protected.HandleFunc("/applications/{id}/review", appHandler.AddReview).Methods("POST")

	admin := api.PathPrefix("/admin").Subrouter()
	admin.Use(authHandler.AdminMiddleware)
	admin.HandleFunc("/applications", appHandler.GetAllApplications).Methods("GET")
	admin.HandleFunc("/applications/{id}/status", appHandler.UpdateStatus).Methods("POST")

	staticDir := "./frontend/build"
	if _, err := os.Stat(staticDir); err == nil {
		r.PathPrefix("/").Handler(http.FileServer(http.Dir(staticDir)))
		log.Println("Serving static files from", staticDir)
	} else {
		log.Println("Frontend build not found, API only mode")
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

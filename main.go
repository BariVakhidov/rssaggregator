package main

import (
	"database/sql"
	"fmt"
	"github.com/BariVakhidov/rssaggregator/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln(err)
	}

	portString := os.Getenv("PORT")

	if portString == "" {
		log.Fatal("No PORT provided")
	}

	dbPortString := os.Getenv("DB_URL")

	if dbPortString == "" {
		log.Fatal("No DB_URL provided")
	}

	dbConnection, dbErr := sql.Open("postgres", dbPortString)
	if dbErr != nil {
		log.Fatalln("Can't connect to db: ", dbErr)
	}

	apiCfg := apiConfig{DB: database.New(dbConnection)}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowCredentials: false,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodOptions},
		AllowedOrigins:   []string{"https://*", "http://*"},
		ExposedHeaders:   []string{"Link"},
		MaxAge:           300,
	}))
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	v1Router := chi.NewRouter()
	v1Router.Get("/health", handlerReadiness)
	v1Router.Get("/err", handlerErr)
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUserByAPIKey))

	router.Mount("/v1", v1Router)

	server := &http.Server{Addr: ":" + portString, Handler: router}

	fmt.Printf("Server starting on port %v", server.Addr)
	serverErr := server.ListenAndServe()

	if serverErr != nil {
		log.Fatalln("Unable to start server: ", serverErr)
	}
}

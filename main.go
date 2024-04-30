package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln(err)
	}

	portString := os.Getenv("PORT")

	if portString == "" {
		log.Fatal("No PORT provided")
	}

	fmt.Println("PORT is ", portString)

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
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)

	router.Mount("/v1", v1Router)

	server := &http.Server{Addr: ":" + portString, Handler: router}

	fmt.Printf("Server starting on port %v", server.Addr)
	serverErr := server.ListenAndServe()

	if serverErr != nil {
		log.Fatalln("Unable to start server: ", serverErr)
	}
}

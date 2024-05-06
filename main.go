package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/BariVakhidov/rssaggregator/api/handlers"
	"github.com/BariVakhidov/rssaggregator/api/routes"
	"github.com/BariVakhidov/rssaggregator/db"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}

	portString := os.Getenv("PORT")

	if portString == "" {
		log.Fatal("No PORT provided")
	}

	dbInstance, dbErr := db.ConnectToDB()
	if dbErr != nil {
		log.Fatalln(dbErr)
	}

	apiCfg := &handlers.ApiConfig{DB: dbInstance}

	go startScrapping(dbInstance, 10, time.Minute)

	router := routes.InitRouter(apiCfg)

	server := &http.Server{Addr: ":" + portString, Handler: router}

	log.Printf("Server starting on port %v", server.Addr)
	serverErr := server.ListenAndServe()

	if serverErr != nil {
		log.Fatalln("Unable to start server: ", serverErr)
	}
}

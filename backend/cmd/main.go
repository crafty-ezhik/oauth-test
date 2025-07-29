package main

import (
	"fmt"
	"github.com/crafty-ezhik/oauth-test/handlers"
	"github.com/crafty-ezhik/oauth-test/routes"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"log/slog"
	"net/http"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error(err.Error())
	}

	router := chi.NewRouter()

	authHandler := handlers.NewAuthHandler()

	routes.InitMiddleware(router, time.Second*10)
	routes.InitRoutes(router, authHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("Listening on port 8080")
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		panic(err)
	}
}

package generator

import "fmt"

func generateMain() string {
	return fmt.Sprintf(`package main

import (
"backforge/internal/app"
	"backforge/internal/database"
	"backforge/internal/routes"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	db := database.InitDB()
	container := app.NewConatiner(db)
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	routes.RegisterRoutes(router, container)

	server := http.Server{
		Handler: router,
		Addr:    ":8080",
	}

	fmt.Println("Server running on :8080")

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
`)
}

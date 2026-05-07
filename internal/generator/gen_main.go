package generator

import (
	"fmt"

	"github.com/eslamward/backforge/internal/parser"
)

func generateMain(cfg *parser.ServerConfig) string {
	dbPort := cfg.Port
	dbEnv := cfg.Env

	if dbPort == "" {
		dbPort = "8080"
	}
	if dbEnv == "" {
		dbEnv = "dev"
	}

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
		Addr:    ":%s",
	}

	fmt.Println("Server running on :%s on %s mode")

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
`, dbPort, dbPort, dbEnv)
}

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/CaioTeixeira95/logic-exp/migrations"
	"github.com/CaioTeixeira95/logic-exp/pkg/app"
	"github.com/CaioTeixeira95/logic-exp/pkg/db"
	"github.com/CaioTeixeira95/logic-exp/pkg/handlers"
	"github.com/CaioTeixeira95/logic-exp/pkg/repositories"
	"github.com/CaioTeixeira95/logic-exp/pkg/services"
	"github.com/gin-gonic/gin"
	migrate "github.com/rubenv/sql-migrate"
)

func main() {
	r := gin.Default()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL can't be empty")
	}

	conn, err := db.Open(dbURL)
	if err != nil {
		log.Fatalf("error connecting to the database: %s", err.Error())
	}

	_, err = db.Migrate(dbURL, migrate.Up, 0, http.FS(migrations.FS))
	if err != nil {
		log.Fatalf("error applying migrations to the database: %s", err.Error())
	}

	repository := repositories.NewRepository(repositories.WithDatabaseOption(conn))

	// Services
	expressionService := services.NewExpressionService(
		services.WithExpressionRepositoryOption(repository),
	)

	// Handlers
	expressionHandler := handlers.NewExpressionHandler(
		handlers.WithExpressionServiceOption(expressionService),
	)

	s := app.NewServer(
		r,
		app.WithExpressionHandlerOption(expressionHandler),
	)

	if err := s.Run(); err != nil {
		log.Fatalf("error running server: %s", err.Error())
	}
}

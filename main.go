package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/repository/postgres"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/repository/postgres/user"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/role"
	router "github.com/ITA-Dnipro/Dp-210_Go/internal/server/http"
	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// Main function
func main() {
	fmt.Println("Start webapp dp210go OK!")
	logger, _ := zap.NewProduction()
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&timezone=utc",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("POSTGRES_DB"))
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(fmt.Errorf("creating db: %w", err))
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		log.Fatal(fmt.Errorf("ping db %s : %w", dsn, err))
	}

	migrationsPath := "migrations"
	err = postgres.MigrateUp(migrationsPath, dsn)
	if err != nil {
		log.Fatal(fmt.Errorf("db migrations: %w", err))
	}
	// TODO remove. for testing purpose.
	repo := user.NewRepository(db)
	hash, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.MinCost)
	repo.Create(context.Background(), entity.User{
		ID:             "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		Name:           "admin",
		Email:          "admin@admin.com",
		PasswordHash:   hash,
		PermissionRole: role.Admin,
	})
	hash, _ = bcrypt.GenerateFromPassword([]byte("operator"), bcrypt.MinCost)
	repo.Create(context.Background(), entity.User{
		ID:             "e4044a74-6557-4c3b-b2d8-4ef933430cf9",
		Name:           "operator",
		Email:          "operator@admin.com",
		PasswordHash:   hash,
		PermissionRole: role.Operator,
	})
	hash, _ = bcrypt.GenerateFromPassword([]byte("user"), bcrypt.MinCost)
	repo.Create(context.Background(), entity.User{
		ID:             "35ce783d-7f09-4ef1-bc27-8bddf1be24d3",
		Name:           "test",
		Email:          "test@admin.com",
		PasswordHash:   hash,
		PermissionRole: role.Viewer,
	})

	logger.Info("starting web server")
	r := router.NewRouter(db, logger)
	// Start server
	log.Fatal(http.ListenAndServe(":8000", r))
}

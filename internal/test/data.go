package test

import (
	"context"
	"database/sql"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/repository/postgres/user"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/role"
	"golang.org/x/crypto/bcrypt"
)

func InitTestData(db *sql.DB) {
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
}

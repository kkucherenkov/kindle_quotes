package users

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/shaj13/go-guardian/v2/auth"
)

type UsersDAO interface {
	ValidateUser(ctx context.Context, r *http.Request, userName, password string) (auth.Info, error)
}

type UsersRepsitory struct {
	db *pgxpool.Pool
}

func NewUserRepository(d *pgxpool.Pool) UsersDAO {
	return &UsersRepsitory{db: d}
}

func (u *UsersRepsitory) ValidateUser(ctx context.Context, r *http.Request, userName, password string) (auth.Info, error) {
	// here connect to db or any other service to fetch user and validate it.
	if userName == "admin" && password == "admin" {
		return auth.NewDefaultUser("admin", "1", nil, nil), nil
	}

	return nil, fmt.Errorf("invalid credentials")
}

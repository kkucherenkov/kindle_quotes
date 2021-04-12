package users

import (
	"context"
	"fmt"
	"net/http"

	"crypto/sha256"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/shaj13/go-guardian/v2/auth"
	"golang.org/x/crypto/pbkdf2"
)

const (
	salt = "kindle_quotes"

	sqlValidateUser = `
SELECT id FROM tbl_users  WHERE username = $1 and password = $2
	`
	sqlCreateUser = `
INSERT INTO tbl_users (username, password)
VALUES ($1, $2)
	`
)

type UsersDAO interface {
	ValidateUser(ctx context.Context, r *http.Request, userName, password string) (auth.Info, error)
	CreateUser(ctx context.Context, r *http.Request, userName, password string) (auth.Info, error)
}

type UsersRepsitory struct {
	db *pgxpool.Pool
}

func NewUserRepository(d *pgxpool.Pool) UsersDAO {
	return &UsersRepsitory{db: d}
}

func (u *UsersRepsitory) CreateUser(ctx context.Context, r *http.Request, userName, password string) (auth.Info, error) {

	conn, err := u.db.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't acuire connection")
	}
	tr, err := conn.Begin(ctx)
	rr, err := tr.Query(ctx, sqlCreateUser, userName, hashPassword(password, salt))

	if err != nil {
		return nil, err
	}
	rr.Close()
	tr.Commit(ctx)
	conn.Release()

	user_id := -1
	conn, err = u.db.Acquire(ctx)
	err = conn.QueryRow(ctx, sqlValidateUser, userName, hashPassword(password, salt)).Scan(&user_id)

	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	conn.Release()
	if user_id >= 0 {
		return auth.NewDefaultUser(userName, fmt.Sprint(user_id), nil, nil), nil
	}
	return nil, fmt.Errorf("invalid credentials")

}

func (u *UsersRepsitory) ValidateUser(ctx context.Context, r *http.Request, userName, password string) (auth.Info, error) {

	user_id := -1
	conn, err := u.db.Acquire(ctx)
	err = conn.QueryRow(ctx, sqlValidateUser, userName, hashPassword(password, salt)).Scan(&user_id)

	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	conn.Release()
	if user_id >= 0 {
		return auth.NewDefaultUser(userName, fmt.Sprint(user_id), nil, nil), nil
	}

	return nil, fmt.Errorf("invalid credentials")
}

func hashPassword(passwd, salt string) string {
	tempPasswd := pbkdf2.Key([]byte(passwd), []byte(salt), 4096, 32, sha256.New)
	return fmt.Sprintf("%x", tempPasswd)
}

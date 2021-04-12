package transport

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kkucherenkov/kindle_quotes/pkg/quotes"

	"github.com/shaj13/go-guardian/v2/auth"
	"github.com/shaj13/go-guardian/v2/auth/strategies/basic"
	"github.com/shaj13/go-guardian/v2/auth/strategies/jwt"
	"github.com/shaj13/go-guardian/v2/auth/strategies/union"
	"github.com/shaj13/libcache"
)

type HttpHandler interface {
	GetQuotes() http.HandlerFunc
	GetBooks() http.HandlerFunc
	GetAuthors() http.HandlerFunc
	GetQuotesByAuthor() http.HandlerFunc
	GetQuotesByTitle() http.HandlerFunc

	Login() http.HandlerFunc
	Registration() http.HandlerFunc

	// GetQuotes(w http.ResponseWriter, r *http.Request)
	// GetBooks(w http.ResponseWriter, r *http.Request)
	// GetAuthors(w http.ResponseWriter, r *http.Request)
	// GetQuotesByAuthor(w http.ResponseWriter, r *http.Request)
	// GetQuotesByTitle(w http.ResponseWriter, r *http.Request)

	// Login(w http.ResponseWriter, r *http.Request)
	// Registration(w http.ResponseWriter, r *http.Request)
}

type KQHandler struct {
	db          *pgxpool.Pool
	qRepository *quotes.DBQuotesRepository
	strategy    union.Union
	keeper      jwt.SecretsKeeper
}

func New(d *pgxpool.Pool, repo *quotes.DBQuotesRepository) HttpHandler {
	keeper := jwt.StaticSecret{
		ID:        "secret-id",
		Secret:    []byte("secret"),
		Algorithm: jwt.HS256,
	}
	cache := libcache.FIFO.New(0)
	cache.SetTTL(time.Minute * 5)
	cache.RegisterOnExpired(func(key, _ interface{}) {
		cache.Peek(key)
	})
	basicStrategy := basic.NewCached(validateUser, cache)
	jwtStrategy := jwt.New(cache, keeper)
	strategy := union.New(jwtStrategy, basicStrategy)
	return KQHandler{db: d, qRepository: repo, keeper: keeper, strategy: strategy}
}

func (h KQHandler) GetQuotes() http.HandlerFunc
func (h KQHandler) GetBooks() http.HandlerFunc
func (h KQHandler) GetAuthors() http.HandlerFunc
func (h KQHandler) GetQuotesByAuthor() http.HandlerFunc
func (h KQHandler) GetQuotesByTitle() http.HandlerFunc

func (h KQHandler) Login() http.HandlerFunc {
	return h.middleware(http.HandlerFunc(h.login))
}
func (h KQHandler) Registration() http.HandlerFunc

func (h KQHandler) login(w http.ResponseWriter, r *http.Request) {
	u := auth.User(r)
	token, _ := jwt.IssueAccessToken(u, h.keeper)
	body := fmt.Sprintf("token: %s \n", token)
	w.Write([]byte(body))
}

func (h KQHandler) registration(w http.ResponseWriter, r *http.Request) {

}

func (h KQHandler) middleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing Auth Middleware")
		_, user, err := h.strategy.AuthenticateRequest(r)
		if err != nil {
			fmt.Println(err)
			code := http.StatusUnauthorized
			http.Error(w, http.StatusText(code), code)
			return
		}
		log.Printf("User %s Authenticated\n", user.GetUserName())
		r = auth.RequestWithUser(user, r)
		next.ServeHTTP(w, r)
	})
}

func validateUser(ctx context.Context, r *http.Request, userName, password string) (auth.Info, error) {
	// here connect to db or any other service to fetch user and validate it.
	if userName == "admin" && password == "admin" {
		return auth.NewDefaultUser("admin", "1", nil, nil), nil
	}

	return nil, fmt.Errorf("Invalid credentials")
}

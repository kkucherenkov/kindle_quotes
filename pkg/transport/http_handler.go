package transport

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kkucherenkov/kindle_quotes/pkg/quotes"
	"github.com/kkucherenkov/kindle_quotes/pkg/users"

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
	qRepository quotes.DBQuotesRepository
	strategy    union.Union
	keeper      jwt.SecretsKeeper
}

func New(qr quotes.DBQuotesRepository, ur users.UsersDAO) HttpHandler {
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

	basicStrategy := basic.NewCached(ur.ValidateUser, cache)
	jwtStrategy := jwt.New(cache, keeper)
	strategy := union.New(jwtStrategy, basicStrategy)
	return KQHandler{qRepository: qr, keeper: keeper, strategy: strategy}
}

func (h KQHandler) GetQuotes() http.HandlerFunc {
	return nil
}
func (h KQHandler) GetBooks() http.HandlerFunc {
	return nil
}
func (h KQHandler) GetAuthors() http.HandlerFunc {
	return nil
}
func (h KQHandler) GetQuotesByAuthor() http.HandlerFunc {
	return nil
}
func (h KQHandler) GetQuotesByTitle() http.HandlerFunc {
	return nil
}

func (h KQHandler) Login() http.HandlerFunc {
	return h.middleware(http.HandlerFunc(h.login))
}
func (h KQHandler) Registration() http.HandlerFunc {
	return http.HandlerFunc(h.registration)
}

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

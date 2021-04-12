package transport

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
}

type KQHandler struct {
	qRepository quotes.DBQuotesRepository
	userDAO     users.UsersDAO
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
	return KQHandler{qRepository: qr, keeper: keeper, strategy: strategy, userDAO: ur}
}

func (h KQHandler) GetQuotes() http.HandlerFunc {
	return h.middleware(http.HandlerFunc(h.getQuotes))
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

func (h KQHandler) getQuotes(w http.ResponseWriter, r *http.Request) {
	body := fmt.Sprintf("status: %s \n", "It works")
	w.Write([]byte(body))
}

func (h KQHandler) registration(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	data := map[string]string{}
	json.Unmarshal(body, &data)
	// fmt.Println(data)
	auth, err := h.userDAO.CreateUser(r.Context(), r, data["username"], data["password"])

	token, _ := jwt.IssueAccessToken(auth, h.keeper)
	newBody := fmt.Sprintf("token: %s \n", token)
	w.Write([]byte(newBody))
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

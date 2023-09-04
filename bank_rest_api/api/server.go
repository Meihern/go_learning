package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/meihern/go_learning/storage"
)

func writeJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func createJWT(subject string) (string, error) {
	now := time.Now()
	claims := &jwt.RegisteredClaims{
		ID:        uuid.NewString(),
		IssuedAt:  &jwt.NumericDate{Time: now},
		ExpiresAt: &jwt.NumericDate{Time: now.Add(time.Duration(5) * time.Minute)},
		NotBefore: &jwt.NumericDate{Time: now},
		Subject:   subject,
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func withJWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("x-jwt-token")

		_, err := validateJWT(tokenString)
		if err != nil {
			writeJson(w, http.StatusForbidden, ApiError{Error: "Invalid token"})
			return
		}

		handlerFunc(w, r)
	}

}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method : %v", token.Header["alg"])
		}

		if err := token.Claims.Valid(); err != nil {
			return nil, err
		}

		return []byte(secret), nil
	})
}

func getIDFromRequest(r *http.Request) (uuid.UUID, error) {
	id := mux.Vars(r)["id"]

	if uuid, err := uuid.Parse(id); err != nil {
		return uuid, fmt.Errorf("Invalid ID: %s", id)
	} else {
		return uuid, nil
	}

}

type ApiError struct {
	Error     string    `json:"error"`
	TimeStamp time.Time `json:"timestamp"`
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			writeJson(w, http.StatusBadRequest, ApiError{Error: err.Error(), TimeStamp: time.Now().In(time.UTC)})
		}
	}
}

type APIServer struct {
	listenAddr string
	store      storage.Storage
}

func NewAPIServer(listenAddr string, store storage.Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	s.registerHandlerAccounts(router)

	log.Println("Bank Rest Api Server running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

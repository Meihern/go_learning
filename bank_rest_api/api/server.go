package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_"github.com/lib/pq"
	"github.com/meihern/go_learning/storage"
)

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func GetIDFromRequest(r *http.Request) (uuid.UUID, error) {
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

func MakeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJson(w, http.StatusBadRequest, ApiError{Error: err.Error(), TimeStamp: time.Now().In(time.UTC)})
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

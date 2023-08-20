package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *APIServer) registerHandlerAccounts(router *mux.Router) {

	router.HandleFunc("/accounts", MakeHTTPHandleFunc(s.handleAccounts))
	router.HandleFunc("/accounts/{id}", MakeHTTPHandleFunc(s.handleAccountsWithID))
	router.HandleFunc("/accounts/{id}/transfer", MakeHTTPHandleFunc(s.handleTransfer))

}

func (s *APIServer) handleAccounts(w http.ResponseWriter, r *http.Request) error {

	if r.Method == "GET" {
		return s.handleGetAccounts(w, r)
	}

	if r.Method ==  "POST" {
		return s.handleCreateAccount(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)

}

func (s *APIServer) handleAccountsWithID(w http.ResponseWriter, r *http.Request) error {

	switch m := r.Method; m {
	case "GET":
		return s.handleGetAccount(w, r)
	case "PUT":
		return s.handleUpdateAccount(w, r)
	case "DELETE":
		return s.handleDeleteAccount(w, r)
	default:
		 return fmt.Errorf("method not allowed %s", m)
	}

}

func (s *APIServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {
	return nil

}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}

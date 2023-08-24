package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/meihern/go_learning/types"
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

	if r.Method == "POST" {
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
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, accounts)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := new(types.CreateOrUpdateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return err
	}

	account := types.NewAccount(createAccountReq.FirstName, createAccountReq.LastName)
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJson(w, http.StatusCreated, account)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := GetIDFromRequest(r)
	if err != nil {
		return err
	}

	account, err := s.store.GetAccountByID(id)
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, account)
}

func (s *APIServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := GetIDFromRequest(r)
	if err != nil {
		return err
	}

	account, err := s.store.GetAccountByID(id)

	updateAccountReq := new(types.CreateOrUpdateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(updateAccountReq); err != nil {
		return err
	}

	account.FirstName = updateAccountReq.FirstName
	account.LastName = updateAccountReq.LastName

	if err := s.store.UpdateAccount(account); err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, account)

}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := GetIDFromRequest(r)
	if err != nil {
		return err
	}

	return s.store.DeleteAccount(id)
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	return nil
}

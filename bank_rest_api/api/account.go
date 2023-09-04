package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/meihern/go_learning/types"
)

func (s *APIServer) registerHandlerAccounts(router *mux.Router) {

	router.HandleFunc("/accounts", makeHTTPHandleFunc(s.handleAccounts))
	router.HandleFunc("/accounts/{id}", withJWTAuth(makeHTTPHandleFunc(s.handleAccountsWithID)))
	router.HandleFunc("/accounts/{id}/transfer", makeHTTPHandleFunc(s.handleTransfer))

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

	return writeJson(w, http.StatusOK, accounts)
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

	token, err := createJWT(account.ID.String())

	if err != nil {
		return err
	}

	fmt.Println("JWT TOKEN : ", token)

	return writeJson(w, http.StatusCreated, account)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getIDFromRequest(r)
	if err != nil {
		return err
	}

	account, err := s.store.GetAccountByID(id)
	if err != nil {
		return err
	}

	return writeJson(w, http.StatusOK, account)
}

func (s *APIServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getIDFromRequest(r)
	if err != nil {
		return err
	}

	account, err := s.store.GetAccountByID(id)
	if err != nil {
		return err
	}

	updateAccountReq := new(types.CreateOrUpdateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(updateAccountReq); err != nil {
		return err
	}

	account.FirstName = updateAccountReq.FirstName
	account.LastName = updateAccountReq.LastName
	account.UpdatedAt = time.Now().In(time.UTC)

	if err := s.store.UpdateAccount(account); err != nil {
		return err
	}

	return writeJson(w, http.StatusOK, account)

}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getIDFromRequest(r)
	if err != nil {
		return err
	}

	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}

	return writeJson(w, http.StatusOK, map[string]string{"message": "Deleted Account: " + id.String()})

}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	return nil
}

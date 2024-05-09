package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// APIServer represents the JSON API server.
type APIServer struct {
	listenAddr string
	store      Storage
}

// newAPIServer creates a new instance of APIServer with the specified listen address and storage interface.
func newAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

// Run starts the JSON API server.
func (s *APIServer) Run() {
	router := mux.NewRouter()

	// Define routes and link them to corresponding handler functions.
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleGetAccountByID))
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleCreateAccount))
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleDeleteAccount))
	router.HandleFunc("/transfer", makeHTTPHandleFunc(s.handleTransfer))

	log.Println("JSON API server running on port: ", s.listenAddr)

	//start the server
	http.ListenAndServe(s.listenAddr, router)
}

// handleAccount handles requests to /account endpoint.
func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	// Handle different HTTP methods for /account endpoint.
	switch r.Method {
	case "GET":
		return s.handleGetAccount(w, r)
	case "POST":
		return s.handleCreateAccount(w, r)
	default:
		return fmt.Errorf("method not allowed: %s", r.Method)
	}

}

// handleGetAccount handles GET requests to /account endpoint.
func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)

}

// handleGetAccountByID handles GET requests to /account/{id} endpoint.
func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	//handle GET request to get specific acc info
	if r.Method == "GET" {
		id, err := getID(r)
		if err != nil {
			return err
		}

		account, err := s.store.GetAccountByID(id)
		if err != nil {
			return err
		}

		return WriteJSON(w, http.StatusOK, account)

	}

	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

// handleCreateAccount handles POST requests to /account endpoint.
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	//handle POST request to create new acc
	CreateAccountReq := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(CreateAccountReq); err != nil {
		return err
	}

	account := NewAccount(CreateAccountReq.FirstName, CreateAccountReq.LastName)
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

// handleDeleteAccount handles DELETE requests to /account endpoint.
func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	//handle DELETE request to delete an acc
	id, err := getID(r)
	if err != nil {
		return err
	}

	if err := s.store.DeleteAccount((id)); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

// handleTransferAccount handles POST requests to /account/transfer endpoint.
func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	//handle POST request to transfer func btw accounts
	transferRequest := new(TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferRequest); err != nil {
		return err
	}
	defer r.Body.Close()
	return WriteJSON(w, http.StatusOK, transferRequest)
}

// writes JSON response to the http.ResponseWriter with the specified status code.
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// apiFunc is a function signature for handler functions.
type apiFunc func(http.ResponseWriter, *http.Request) error

// ApiError represents the structure for API error responses.
type ApiError struct {
	Error string `json:"error"`
}

// makeHTTPHandleFunc is a decorator to convert an apiFunc to an http.HandlerFunc
func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}
	return id, nil

}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

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
	Error string
}

// makeHTTPHandleFunc is a decorator to convert an apiFunc to an http.HandlerFunc
func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

// APIServer represents the JSON API server.
type APIServer struct {
	listenAddr string
}

// newAPIServer creates a new instance of APIServer with the specified listen address.
func newAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

// Run starts the JSON API server.
func (s *APIServer) Run() {
	router := mux.NewRouter()

	// Define routes and link them to corresponding handler functions.
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleGetAccount))
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleCreateAccount))
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleDeleteAccount))
	router.HandleFunc("/account/transfer", makeHTTPHandleFunc(s.handleTransferAccount))

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
	case "DELETE":
		return s.handleDeleteAccount(w, r)
	default:
		return fmt.Errorf("method not allowed: %s", r.Method)
	}

}

// handleGetAccount handles GET requests to /account endpoint.
func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	//handle GET request to get specific acc info
	id := mux.Vars(r)["id"]
	fmt.Println(id)
	return WriteJSON(w, http.StatusOK, &Account{})
}

// handleCreateAccount handles POST requests to /account endpoint.
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	//handle POST request to create new acc
	return nil
}

// handleDeleteAccount handles DELETE requests to /account endpoint.
func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	//handle DELETE request to delete an acc
	return nil
}

// handleTransferAccount handles POST requests to /account/transfer endpoint.
func (s *APIServer) handleTransferAccount(w http.ResponseWriter, r *http.Request) error {
	//handle POST request to transfer func btw accounts
	return nil
}

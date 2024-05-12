package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
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
	router.HandleFunc("/login", makeHTTPHandleFunc(s.handleLogin))
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", WithJWTAuth(makeHTTPHandleFunc(s.handleGetAccountByID), s.store))
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleCreateAccount))
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleDeleteAccount))
	router.HandleFunc("/transfer", makeHTTPHandleFunc(s.handleTransfer))

	log.Println("JSON API server running on port: ", s.listenAddr)

	//start the server
	http.ListenAndServe(s.listenAddr, router)
}


func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	// Only allow POST requests for authentication.
	if r.Method != "POST" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	// Retrieve account information from the store.
	acc, err := s.store.GetAccountByNumber(int(req.Number))
	if err != nil {
		return err
	}

	// Validate the provided password.
	if !acc.ValidPassword(req.Password) {
		return fmt.Errorf("not authenticated")
	}

	// Create JWT token for authentication.
	token, err := createJWT(acc)
	if err != nil {
		return err
	}

	// Prepare and send the login response.
	resp := LoginResponse{
		Token:  token,
		Number: acc.Number,
	}

	return WriteJSON(w, http.StatusOK, resp)
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
func (s *APIServer) handleGetAccount(w http.ResponseWriter, _ *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}

	// Return the accounts as JSON response.
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

		// Retrieve account by ID from the store.
		account, err := s.store.GetAccountByID(id)
		if err != nil {
			return err
		}

		return WriteJSON(w, http.StatusOK, account)

	}

	// Handle DELETE request to delete an account.
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

// handleCreateAccount handles POST requests to /account endpoint.
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	//handle POST request to create new acc
	createAccountReq := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return err
	}

	// Create a new account.
	account, err := NewAccount(createAccountReq.FirstName, createAccountReq.LastName, createAccountReq.Password)
	if err != nil {
		return err
	}

	// Store the new account.
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	 // Return the new account as JSON response.
	return WriteJSON(w, http.StatusOK, account)
}

// handleDeleteAccount handles DELETE requests to /account endpoint.
func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	//handle DELETE request to delete an acc
	id, err := getID(r)
	if err != nil {
		return err
	}

	 // Delete the account from the store.
	if err := s.store.DeleteAccount((id)); err != nil {
		return err
	}

	// Return the ID of the deleted account as JSON response.
	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

// handleTransfer handles POST requests to /account/transfer endpoint.
func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	//handle POST request to transfer func btw accounts
	transferRequest := new(TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferRequest); err != nil {
		return err
	}
	defer r.Body.Close()

	// Return the transfer request as JSON response.
	return WriteJSON(w, http.StatusOK, transferRequest)
}

// writes JSON response to the http.ResponseWriter with the specified status code.
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

// createJWT creates a JWT token for the given account.
func createJWT(account *Account) (string, error) {
	// Create the Claims
	claims := &jwt.MapClaims{
		"expiresAt":     jwt.NewNumericDate(time.Unix(1516239022, 0)),
		"accountNumber": account.Number,
	}

	 // Get the secret key from the environment variables
	secret := os.Getenv("JWT_SECRET")

	 // Create a new JWT token with the claims and signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	 // Sign the token with the secret key and return the signed token as a string
	return token.SignedString([]byte(secret))
}

// permissionDenied sends a permission denied response with the appropriate HTTP status code.
func permissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, ApiError{Error: "permission denied"})
}


// WithJWTAuth is a middleware function that authenticates requests using JWT tokens.
// It takes an HTTP handler function and a storage interface as input, and returns
// a new HTTP handler function. This middleware verifies the JWT token in the request header
// and authorizes access based on the token's validity and the account associated with it.
func WithJWTAuth(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("calling JWT auth middleware")

		// Retrieve the JWT token from the request header
		tokenString := r.Header.Get("x-jwt-token")
		// Validate the JWT token
		token, err := validateJWT(tokenString)
		if err != nil {
			permissionDenied(w)
			return
		}

		if !token.Valid {
			permissionDenied(w)
			return
		}

		// Retrieve the user ID from the request
		userID, err := getID(r)
		if err != nil {
			permissionDenied(w)
			return
		}

		 // Retrieve the account associated with the user ID from the storage
		account, err := s.GetAccountByID(userID)
		if err != nil {
			permissionDenied(w)
			return
		}

		// Extract the claims from the JWT token
		claims := token.Claims.(jwt.MapClaims)
	
		// Compare the account number in the claims with the account number associated with the user
		if account.Number != int64(claims["accountNumber"].(float64)) {
			permissionDenied(w)
			return
		}

		if err != nil {
			WriteJSON(w, http.StatusForbidden, ApiError{Error: "invalid token"})
			return
		}

		// If all checks pass, call the original handler function
		handlerFunc(w, r)

	}
}

// validateJWT validates the JWT token. It takes the JWT token string as input
// and returns the parsed token along with any error encountered during validation.
func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	// Parse and validate the JWT token using the secret key
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
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
		// Call the provided handler function and handle any errors
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

// getID extracts and parses the 'id' parameter from the request URL.
func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}
	return id, nil
}

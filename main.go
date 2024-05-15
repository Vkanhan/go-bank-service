package main

import (
	"flag"
	"fmt"
	"log"
)

// seedAccount creates a new account using the provided store and user details, then prints the account number.
func seedAccount(store Storage, firstName, lastName, password string) *Account {
	// Create a new account with the given first name, last name, and password.
	account, err := NewAccount(firstName, lastName, password)
	if err != nil {
		log.Fatal(err)
	}

	// Store the created account in the database.
	if err := store.CreateAccount(account); err != nil {
		log.Fatal(err)
	}

	// Print the new account number.
	fmt.Println("new account => ", account.Number)

	// Return the created account.
	return account

}

// seedAccounts seeds the database with predefined accounts.
func seedAccounts(s Storage) {
	seedAccount(s, "albert", "einstein", "travel")
}

func main() {

	// Define a command-line flag for seeding the database.
	seed := flag.Bool("seed", false, "seed the db")
	flag.Parse()

	// Initialize a new instance of PostgresStore.
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the database by creating necessary tables.
	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	// Check if the seed flag is set to true.
	if *seed {
		fmt.Println("seeding the database")

		//seeed stuff
		seedAccounts(store)
	}

	// Create a new instance of APIServer listening on port 3000.
	// and pass the initialized store as a parameter.
	server := newAPIServer(":3000", store)

	// Run the server to start handling incoming HTTP requests.
	server.Run()

}

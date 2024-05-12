package main

import (
	"flag"
	"fmt"
	"log"
)

func seedAccount(store Storage, firstName, lastName, password string) *Account {
	account, err := NewAccount(firstName, lastName, password)
	if err != nil {
		log.Fatal(err)
	}

	if err := store.CreateAccount(account); err != nil {
		log.Fatal(err)
	}

	fmt.Println("new account => ", account.Number)

	return account

}

func seedAccounts(s Storage) {
	seedAccount(s, "albert", "einstein", "travel")
}

func main() {

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

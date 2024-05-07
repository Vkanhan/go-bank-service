package main

import (
	"log"
)

func main() {

	// Initialize a new instance of PostgresStore.
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the database by creating necessary tables.
	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	// Create a new instance of APIServer listening on port 3000.
	// and pass the initialized store as a parameter.
	server := newAPIServer(":3000", store)

	// Run the server to start handling incoming HTTP requests.
	server.Run()

}

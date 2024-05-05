package main

func main() {

	// Create a new instance of APIServer listening on port 3000.
	server := newAPIServer(":3000")

	// Run the server to start handling incoming HTTP requests.
	server.Run()

}

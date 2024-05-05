package main

func main() {

	server := newAPIServer(":3000")
	server.Run()

}

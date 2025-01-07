package main

import (
	"log"
	"net/http"

	"hrapplication/internal/azure"
	"hrapplication/internal/handlers"
	"hrapplication/internal/utils"
)

func main() {
	// Load environment variables
	env, err := utils.LoadEnvironment()
	if err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	// Retrieve the combined PEM file
	pemFile, err := azure.GetCertAndKey(env)
	if err != nil {
		log.Fatalf("Error retrieving PEM file: %v", err)
	}
	defer utils.RemoveTempFile(pemFile)

	// Fetch the port number
	port, err := azure.GetSecret(env, "SERVERPORT")
	if err != nil || port == "" {
		log.Fatalf("Invalid port number: %v", err)
	}

	// Register handlers
	http.HandleFunc("/", handlers.ServeHomePage)
	http.HandleFunc("/users", handlers.ServeUsersPage)
	http.HandleFunc("/api/users", handlers.HandleUsers)
	http.HandleFunc("/api/users/", handlers.HandleDeleteUser)

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Start the HTTPS server
	log.Printf("Server started on https://localhost:%s", port)
	err = http.ListenAndServeTLS(":443", pemFile, pemFile, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

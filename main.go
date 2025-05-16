package main

import "fmt" 		// Provides formatted I/O functions like Println, Printf, etc.
import "net/http" 	// Provides HTTP client and server implementations

// User struct represents a user with a single field `Name`.
// The `json:"name"` tag specifies how this field will be encoded/decoded in JSON.
type User struct {
	Name string `json:"name"`
}

// userCache is a map that stores users with an integer key.
// This acts as an in-memory storage for user data.
var userCache = make(map[int]User)

// main is the entry point of the program.
func main(){
	// Create a new HTTP request multiplexer (router).
	mux := http.NewServeMux()

	// Register a handler for GET requests to the root URL ("/").
	mux.HandleFunc("/", handleRoot)

    // Register a handler for POST requests to "/users".
    // Note: The route "POST /users" is incorrect and should just be "/users".
	mux.HandleFunc("POST /users", createUser)

    // Print a message to the console indicating the server is starting.
    // Note: There is a typo in `PrintIn`. It should be `Println`.
	fmt.PrintIn("Server listening to :9090")

	// Start the HTTP server on port 9090 and use the `mux` router.
	http.ListenAndServe(":9090", mux)
}

// handleRoot handles GET requests to the root URL ("/").
// It responds with "Hello World".
func handleRoot(
	w http.ResponseWriter, 	// Used to send a response back to the client.
	r *http.Request, 		// Represents the incoming HTTP request.
) {
	// Write "Hello World" to the response.
	fmt.Fprintf(w, "Hello World")
}

// createUser handles POST requests to "/users".
// Currently, this function is empty and does not perform any operations. 
func createUser(
	w http.ResponseWriter,
	r *http.Request,
){
	
}
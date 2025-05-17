package main

import (
	"encoding/json" // Provides functions to encode and decode JSON data.
	"fmt"           // Provides formatted I/O functions like Println, Printf, etc.
	"net/http"      // Provides HTTP client and server implementations.
	"strconv"       // Provides functions to convert strings to numbers and vice versa.
	"sync"          // Provides synchronization primitives such as mutexes for safe concurrent access.
)

// User struct represents a user with a single field `Name`.
// The `json:"name"` tag specifies how this field will be encoded/decoded in JSON.
type User struct {
	Name string `json:"name"`
}

// userCache is a map that stores users with an integer key.
// This acts as an in-memory storage for user data.
var userCache = make(map[int]User)

// cacheMutex is a read-write mutex used to synchronize access to the userCache.
// This ensures thread-safe operations on the map.
var cacheMutex sync.RWMutex

// main is the entry point of the program.
func main() {
	// Create a new HTTP request multiplexer (router).
	mux := http.NewServeMux()

	// Register a handler for GET requests to the root URL ("/").
	mux.HandleFunc("/", handleRoot)

	// Register a handler for POST requests to "/users".
	mux.HandleFunc("POST /users", createUser)

	// Register a handler for GET requests to "/users/{id}".
	mux.HandleFunc("GET /users/{id}", getUser)

	// Register a handler for DELETE requests to "/users/{id}".
	mux.HandleFunc("DELETE /users/{id}", deleteUser)

	// Print a message to the console indicating the server is starting.
	fmt.Println("Server listening to :9090")

	// Start the HTTP server on port 9090 and use the `mux` router.
	http.ListenAndServe(":9090", mux)
}

// handleRoot handles GET requests to the root URL ("/").
// It responds with "Hello World".
func handleRoot(
	w http.ResponseWriter, // Used to send a response back to the client.
	r *http.Request, // Represents the incoming HTTP request.
) {
	// Write "Hello World" to the response.
	fmt.Fprintf(w, "Hello World")
}

// getUser handles GET requests to "/users/{id}".
// It retrieves a user by their ID from the userCache.
func getUser(
	w http.ResponseWriter, // Used to send a response back to the client.
	r *http.Request, // Represents the incoming HTTP request.
) {
	// Extract the user ID from the URL path.
	id, err := strconv.Atoi(r.URL.Path[len("/users/"):])
	if err != nil {
		// If the ID is not a valid integer, return a 400 Bad Request error.
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Lock the cache for reading and retrieve the user.
	cacheMutex.RLock()
	user, ok := userCache[id]
	cacheMutex.RUnlock()

	if !ok {
		// If the user is not found, return a 404 Not Found error.
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	// Convert the user to JSON and send it in the response.
	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(user)
	if err != nil {
		// If there is an error during JSON encoding, return a 500 Internal Server Error.
		// This is a server-side error, so we use http.StatusInternalServerError.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the JSON response with a 200 OK status.
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// createUser handles POST requests to "/users".
// It creates a new user and stores it in the userCache.
func createUser(
	w http.ResponseWriter, // Used to send a response back to the client.
	r *http.Request, // Represents the incoming HTTP request.
) {
	// Decode the JSON body into a User struct.
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		// If the request body is not valid JSON, return a 400 Bad Request error.
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	// Validate that the user's name is not empty.
	if user.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// Lock the cache for writing and add the new user.
	cacheMutex.Lock()
	userCache[len(userCache)+1] = user
	cacheMutex.Unlock()

	// Respond with a 204 No Content status to indicate success.
	w.WriteHeader(http.StatusNoContent)
}

// deleteUser handles DELETE requests to "/users/{id}".
// It deletes a user by their ID from the userCache.
func deleteUser(
	w http.ResponseWriter, // Used to send a response back to the client.
	r *http.Request, // Represents the incoming HTTP request.
) {
	// Extract the user ID from the URL path.
	id, err := strconv.Atoi(r.URL.Path[len("/users/"):])
	if err != nil {
		// If the ID is not a valid integer, return a 400 Bad Request error.
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the user exists in the cache.
	cacheMutex.RLock()
	_, ok := userCache[id]
	cacheMutex.RUnlock()

	if !ok {
		// If the user is not found, return a 404 Not Found error.
		http.Error(w, "user not found", http.StatusBadRequest)
		return
	}

	// Lock the cache for writing and delete the user.
	cacheMutex.Lock()
	delete(userCache, id)
	cacheMutex.Unlock()

	// Respond with a 204 No Content status to indicate success.
	w.WriteHeader(http.StatusNoContent)
}

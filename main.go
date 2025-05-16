package main

import "fmt"
import "net/http"


type User struct {
	Name string `json:"name"`
}

var userCache = make(map[int]User)

func main(){
	mux := http.NewServeMux()
	// get request
	mux.HandleFunc("/", handleRoot)
	// post request
	mux.HandleFunc("POST /users", createUser)

	fmt.PrintIn("Server listening to :9090")
	http.ListenAndServe(":9090", mux)
}

func handleRoot(
	w http.ResponseWriter,
	r *http.Request,
) {
	fmt.Fprintf(w, "Hello World")
}
 
func createUser(
	w http.ResponseWriter,
	r *http.Request,
){
	
}
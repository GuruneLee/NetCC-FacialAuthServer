package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

/*
type Request struct {
	Name  string `json:"name"`
	Image []byte `json:"image"`
}
*/

const (
	SuccesMsg string = "signup success"
	ErrorMsg  string = "you got some errors"
	URL       string = "http://116.89.189.52:8080/get/feature" //face-ai-server's URL
	DB_name   string = "DB.json"
)

/*
// signup/face
func Signup(w http.ResponseWriter, r *http.Request) {

}
*/

//main
func main() {
	fmt.Println("Facial-Auth-server started")
	router := mux.NewRouter()

	// signup/face
	router.HandleFunc("/signup/face", Signup).Methods(http.MethodPost)
	router.HandleFunc("/signin/face", Signin).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":8081", router))
}

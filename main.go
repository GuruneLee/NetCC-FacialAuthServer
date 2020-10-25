package main

import (
	"fmt"
	"log"
	"net/http"

	face "github.com/Kagami/go-face"
	"github.com/gorilla/mux"
)

type Request struct {
	Name  string `json:"name"`
	Image []byte `json:"image"`
}

type Resp struct {
	Feature face.Descriptor `json:"feature"`
	Error   string          `json:"error"`
}

type Meta struct {
	Name string `json:name`
}

const (
	SuccesMsg string = "signup success"
	ErrorMsg  string = "you got some errors"
	URL       string = "http://116.89.189.52:8080/get/feature"
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

	log.Fatal(http.ListenAndServe(":8081", router))
}

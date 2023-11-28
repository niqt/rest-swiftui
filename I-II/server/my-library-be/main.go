package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func main() {
	fmt.Printf("*** Server starting ***\n")
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/books", GetBooksEndpoint).Methods("GET") // ?search=""
	router.HandleFunc("/api/v1/books/{id}", GetBookByIdEndpoint).Methods("GET")
	router.HandleFunc("/api/v1/books/{id}", DeleteBookByIdEndpoint).Methods("DELETE")
	router.HandleFunc("/api/v1/books", SaveBookEndpoint).Methods("POST")
	router.HandleFunc("/api/v1/user", CreateUserEndpoint).Methods("POST")
	router.HandleFunc("/api/v1/user/login", AuthenticateUserEndpoint).Methods("POST")

	userExist, _ := UserExists("x@y.com")
	if !userExist {
		CreateUser(User{Email: "x@y.com", Password: "123456"})
	}

	SaveBook(Book{Title: "Test Book " + time.Now().String(), ISBN13: "12", User: "x@y.com"})

	router.Use(JwtAuthentication)
	log.Fatal(http.ListenAndServe(":8080", router))
}

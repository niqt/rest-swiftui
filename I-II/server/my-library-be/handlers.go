package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

// get the list of the books for the connected user
// use the search string if present
func GetBooksEndpoint(w http.ResponseWriter, req *http.Request) {
	v := req.URL.Query()

	search := v.Get("search")
	req.Context().Value("username")
	json.NewEncoder(w).Encode(GetBooks(search, req.Context().Value("username").(string)))
}

// get the book by id
func GetBookByIdEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id := params["id"]

	json.NewEncoder(w).Encode(GetBookById(id))
}

// delete the book by id
func DeleteBookByIdEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id := params["id"]

	err := DeleteBookById(id)
	if err != nil {
		w.WriteHeader(500) // internal error if it's impossible to delete
	}
	json.NewEncoder(w)
}

// save the book
func SaveBookEndpoint(w http.ResponseWriter, req *http.Request) {
	var book Book

	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = io.ReadAll(req.Body)
	}
	err := json.Unmarshal(bodyBytes, &book)

	if err != nil {
		w.WriteHeader(500) // internal error if it's impossible decode the book
		return
	}
	book.User = req.Context().Value("username").(string)
	result := SaveBook(book)

	json.NewEncoder(w).Encode(result)
}

// create the user
func CreateUserEndpoint(w http.ResponseWriter, req *http.Request) {
	var user User

	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = io.ReadAll(req.Body)
	}
	err := json.Unmarshal(bodyBytes, &user)

	if err != nil {
		w.WriteHeader(500) // internal error if it's impossible decode the user
		return
	}

	result, err := CreateUser(user)

	if err != nil {
		w.WriteHeader(500)
		return
	}

	json.NewEncoder(w).Encode(result)
}

// login
func AuthenticateUserEndpoint(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	account := &User{}

	err := json.NewDecoder(req.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		w.WriteHeader(500)
		return
	}

	resp, err := Login(account.Email, account.Password)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

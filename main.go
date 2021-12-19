package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

//book struct
type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

//author struct
type Author struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

//init book stuct
var books map[string]Book = make(map[string]Book)

//closure function to generate ids for book
func IdGenerator() func() int {
	id := 3
	return func() int {
		id++
		return id
	}
}

var IdGen func() int = IdGenerator()

//handlers
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

//get a single book
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	if book, exist := books[params["id"]]; exist {
		json.NewEncoder(w).Encode(book)
		return
	}
	json.NewEncoder(w).Encode(&Book{})
}

func createBook(w http.ResponseWriter, r *http.Request) {
	setHeaderJson(w)
	book := Book{}
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(IdGen())
	books[book.ID] = book
	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	setHeaderJson(w)
	params := mux.Vars(r)
	book := Book{}
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = params["id"]
	books[params["id"]] = book
	json.NewEncoder(w).Encode(book)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	setHeaderJson(w)
	params := mux.Vars(r)
	delete(books, params["id"])
	json.NewEncoder(w).Encode(books)
}

func main() {
	router := mux.NewRouter()

	//mock data - @todo: implement database
	books["1"] = Book{ID: "1", Isbn: "12345", Title: "Book 1", Author: &Author{FirstName: "Levi", LastName: "Zwannah"}}

	books["2"] = Book{ID: "2", Isbn: "123456", Title: "Book 2", Author: &Author{FirstName: "Levi2", LastName: "Zwannah2"}}

	router.HandleFunc("/api/books", getBooks).Methods("GET")
	router.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/api/books", createBook).Methods("POST")
	router.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	fmt.Println("Server running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func setHeaderJson(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

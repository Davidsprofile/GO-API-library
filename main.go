package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"errors"
)

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"` // turning into json
	Quantity int    `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
	{ID: "3", Title: "War and Peace", Author: "Lev Tolstoy", Quantity: 6},
}

func getBooks(c *gin.Context) { // this function returns json format

	c.IndentedJSON(http.StatusOK, books) // we get nicely formated json. status and telling that we want to send books
}


func bookById(c *gin.Context) {
	id := c.Param("id")  // path parameter 
	book,err := getBookById(id)
	if err != nil {   // if we have error end operation
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not Found"})
		return 
	}
	c.IndentedJSON(http.StatusOK, book)  // if everything ok proceed
}

func checkoutBook (c *gin.Context) {
	id,ok := c.GetQuery("id")
	
	if !ok{
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter"})
		return 
	}

	book,err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not Found"})
		return
	}

	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not Available"})
		return
	}

	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)
}


func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}

	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)
}



func getBookById(id string) (*book, error) {  // error in case if book not found
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}

	return nil, errors.New("book not found")
}


func createBook(c *gin.Context) { // get data of book created in request
	var newBook book
	if err := c.BindJSON(&newBook); err != nil { // bind json to requested book and check for error
		return
	}
	books = append(books, newBook)              // we append book
	c.IndentedJSON(http.StatusCreated, newBook) // return book with status created
}

func main() {
	router := gin.Default()        // router is coming from gin
	router.GET("/books", getBooks) // we tell which route we want to handle /books. when you go to localhost 8080 this will call function
	router.GET("/books/:id", bookById)
	router.POST("/books", createBook)   // Adding information
	router.PATCH("/checkout", checkoutBook)  // updating information
	router.PATCH("/return", returnBook)		// updating information
	router.Run("localhost:8080") // where we want to run our program
}



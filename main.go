package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type book struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Author   string   `json:"author"`
	Category category `json:"category"`
	Location location `json:"location"`
	ISBN     string   `json:"isbn"`
	Holder   user     `json:"user"`
}

type user struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Mail        string `json:"mail"`
	Password    string `json:"password"`
	BooksOnHold []book `json:"book"`
}

type location struct {
	Longitude float64 `json:"id"`
	Latitude  float64 `json:"latitude"`
	Altitude  float64 `json:"altitude"`
}

type category struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

var books = []book{}

var users = []user{
	{ID: "1", Name: "Alex & Bill", Mail: "test@gmail.com", Password: "12345", BooksOnHold: []book{}},
}

func registerBook(c *gin.Context) {
	var newBook book

	if err := c.BindJSON(&newBook); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func findBookById(id string) (*book, error) {
	fmt.Println(id)

	for i, b := range books {
		fmt.Println(b)
		if b.ID == id {

			return &books[i], nil
		}
	}

	fmt.Println("hello it failed")
	return nil, errors.New("book not found")
}

func pickUpABook(c *gin.Context) {
	id := c.Param("bookId")
	book, err := findBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book does not exist"})
		return
	}

	book.Holder = users[0]

	users[0].BooksOnHold = append(users[0].BooksOnHold, *book)

	fmt.Println(users[0])

	c.IndentedJSON(http.StatusOK, nil)
}

func registerUser(c *gin.Context) {
	var newUser user

	if err := c.BindJSON(newUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	users = append(users, newUser)
	c.IndentedJSON(http.StatusCreated, newUser)
}

func leaveABook(c *gin.Context) {
	var updatedBook book
	id := c.Param("bookId")

	book, err := findBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not found"})
	}

	if err := c.BindJSON(updatedBook); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	book.Location = updatedBook.Location
	c.IndentedJSON(http.StatusOK, book)
}

func main() {
	router := gin.Default()
	router.POST("/books", registerBook)
	router.POST("/user", registerUser)
	router.GET("/books", getBooks)
	router.POST("/pickup/:bookId", pickUpABook)
	router.POST("/leave/:bookId", leaveABook)
	router.Run("localhost:8080")
}

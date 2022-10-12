package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Person struct {
	gorm.Model

	Name  string
	Email string `gorm:"typevarchar(100);unique_index"`
	Books []Book
}

type Book struct {
	gorm.Model

	Title      string
	Author     string
	CallNumber int
	PersonID   int
}

var (
	person = &Person{Name: "Jack", Email: "jack@dmail.com"}
	books  = []Book{
		{Title: "Book 1", Author: "Author 1", CallNumber: 1234, PersonID: 1},
		{Title: "Book 2", Author: "Author 2", CallNumber: 5678, PersonID: 1},
	}
)

var db *gorm.DB
var err error

func main() {

	dialect := os.Getenv("DIALECT")
	host := os.Getenv("HOST")
	dbPort := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	dbname := os.Getenv("NAME")
	dbpassword := os.Getenv("PASSWORD")

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, dbname, dbpassword, dbPort)

	db, err = gorm.Open(dialect, dbURI)

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Database connected successfully")
	}

	defer db.Close()

	db.AutoMigrate(&Person{})
	db.AutoMigrate(&Book{})

	// db.Create(person)
	// for idx := range books {
	// 	db.Create(&books[idx])
	// }

	// API routes
	router := gin.Default()
	router.GET("/people", GetPeople)
	router.GET("/person/:id", getPerson)
	router.POST("/person", createPerson)
	router.PUT("/update/person", updatePerson)
	router.DELETE("/delete/person/:id", deletePerson)

	router.Run("localhost:8080")

}

func getPerson(c *gin.Context) {

	id := c.Param("id")
	var person Person
	var books []Book

	db.First(&person, id)
	db.Model(&person).Related(&books)

	person.Books = books

	c.IndentedJSON(http.StatusOK, person)
}

func GetPeople(c *gin.Context) {

	var people []Person
	db.Find(&people)
	c.IndentedJSON(http.StatusOK, person)
}

func createPerson(c *gin.Context) {
	var newPerson Person

	if err := c.BindJSON(&newPerson); err != nil {
		fmt.Println("Decoding err ", err)
		return
	}

	createdPerson := db.Create(&newPerson)
	err = createdPerson.Error
	if err != nil {
		fmt.Println(err)
	}

	c.IndentedJSON(http.StatusCreated, newPerson)

}

func updatePerson(c *gin.Context) {
	var person Person
	var existingPerson Person
	if err := c.BindJSON(&person); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	db.First(&existingPerson, person.ID)

	existingPerson = person
	db.Update(existingPerson)

	c.IndentedJSON(http.StatusOK, existingPerson)
}

func deletePerson(c *gin.Context) {
	id := c.Param("id")
	var person Person
	db.First(&person, id)
	db.Delete(&person)

	c.IndentedJSON(http.StatusOK, person)
}

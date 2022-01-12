package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const (
	apipath = "/api/v1/books"
)

type lib struct {
	dbHost, dbPass, dbName, dbUser string
}

type Book struct {
	Id, Name, Isbn string
}

func main() {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost:3306"
	}

	dbPass := os.Getenv("DB_PASS")
	if dbPass == "" {
		dbPass = "VMware@123"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "dbuser"
	}

	apiPath := os.Getenv("API_PATH")
	if apiPath == "" {
		apiPath = apipath
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "library"
	}
	l := lib{dbHost, dbPass, dbName, dbUser}
	fmt.Println(l)
	r := mux.NewRouter()
	r.HandleFunc(apiPath, l.getBooks).Methods(http.MethodGet)
	r.HandleFunc(apiPath, l.postBooks).Methods(http.MethodPost)
	http.ListenAndServe(":8080", r)

}

func (l lib) postBooks(w http.ResponseWriter, r *http.Request) {
	//convert the http payload into Go book struct
	book := Book{}
	json.NewDecoder(r.Body).Decode(&book)
	//open connection
	db := l.openConnection()

	//Insert query into db
	insertQuery, err := db.Prepare("insert into books values (?, ?, ?)")
	if err != nil {
		log.Fatalf("Error while preparing query %s\n", err.Error())
	}
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Error while beginning tx %s\n", err.Error())
	}
	_, err = tx.Stmt(insertQuery).Exec(book.Id, book.Name, book.Isbn)
	if err != nil {
		log.Fatalf("Error while inserting tx %s\n", err.Error())
	}
	err = tx.Commit()
	if err != nil {
		log.Fatalf("Error while commiting tx %s\n", err.Error())
	}

	//close connection
	l.closeConnection(db)

}

func (l lib) getBooks(w http.ResponseWriter, r *http.Request) {
	//open a connection
	db := l.openConnection()

	//Read all books
	rows, err := db.Query("select * from books")
	if err != nil {
		log.Fatalf("Error while reading books table %s\n", err.Error())
	}
	books := []Book{}
	for rows.Next() {
		var id, name, isbn string
		err := rows.Scan(&id, &name, &isbn)
		if err != nil {
			log.Fatalf("Error while scanning rows %s\n", err.Error())
		}
		aBook := Book{id, name, isbn}
		books = append(books, aBook)
	}
	json.NewEncoder(w).Encode(books)

	//Close the connection
	l.closeConnection(db)
}

func (l lib) openConnection() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s", l.dbUser, l.dbPass, l.dbHost, l.dbName))
	if err != nil {
		log.Fatalf("Error while opening connection %s\n", err.Error())
	}
	return db
}

func (l lib) closeConnection(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Fatalf("Error while closing connection %s\n", err.Error())
	}
}

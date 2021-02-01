package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

type Product struct {
	gorm.Model
	Code  string
	Name  string
	Price uint
}

func main() {
	// ORM stuff..
	db, dbError := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if dbError != nil {
		log.Fatalf("Failed to open DB: %s\n", dbError.Error())
	}

	migrateErr := db.AutoMigrate(&Product{})
	if migrateErr != nil {
		log.Fatalf("Failed to migrate schema: %s\n", migrateErr.Error())
	}

	// Create test item
	db.Create(&Product{
		Code:  "C15",
		Name:  "TestProduct",
		Price: 32,
	})

	// Router setup...
	r := mux.NewRouter()
	r.HandleFunc("/", mainTest)
	r.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		e := json.NewEncoder(w)
		encodeError := e.Encode(map[string]bool{"ok": true})
		if encodeError != nil {
			log.Printf("An error occured during healthcheck: %s\n", encodeError.Error())
		}
	})
	r.HandleFunc("/products/{key}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var product Product
		db.First(&product, "code = ?", vars["key"])
		e := json.NewEncoder(w)
		encodeError := e.Encode(product)
		if encodeError != nil {
			log.Printf("An error occured during encoding: %s\n", encodeError.Error())
		}
	})

	// Server setup...
	srv := &http.Server{
		Addr:         ":8000",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Println("Server is starting...")
	log.Fatal(srv.ListenAndServe())
}

// Function for testing router
func mainTest(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, writeErr := fmt.Fprint(w, "Hi!")
	if writeErr != nil {
		log.Printf("An error occured while executing 'mainTest' func: %s\n", writeErr.Error())
	}
}

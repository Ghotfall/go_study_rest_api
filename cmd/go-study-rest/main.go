package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"go_study_rest_api/pkg/db"
	"go_study_rest_api/pkg/models"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

func main() {
	// Router setup...
	r := mux.NewRouter()
	r.HandleFunc("/students/{name}", getStudent).Methods("GET")
	//r.HandleFunc("/students/{name}", newStudent).Methods("POST")

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

func getStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var student models.Student
	result := db.DB.First(&student, "firstname = ?", vars["name"])

	// For response
	e := json.NewEncoder(w)
	var encodeError error

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		encodeError = e.Encode(map[string]string{"error": "Record not found"})
	} else {
		encodeError = e.Encode(student)
	}

	if encodeError != nil {
		log.Printf("An error occured during encoding: %s\n", encodeError.Error())
	}
}

//func newStudent(w http.ResponseWriter, r *http.Request) {
//
//}

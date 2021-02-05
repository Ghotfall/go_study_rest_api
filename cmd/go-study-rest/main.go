package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"go_study_rest_api/pkg/db"
	"go_study_rest_api/pkg/models"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	// Router setup...
	r := mux.NewRouter()
	r.HandleFunc("/students/{name}", getStudent).Methods("GET")
	r.HandleFunc("/students/{name}", newStudent).Methods("POST")

	// Server setup...
	srv := &http.Server{
		Addr:         ":8000",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	go func() {
		log.Println("Server is starting...")
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	shutdownError := srv.Shutdown(ctx)
	if shutdownError != nil {
		log.Println(shutdownError.Error())
	}
	os.Exit(0)
}

func getStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var student models.Student
	result := db.DB.First(&student, "firstname = ?", vars["name"])

	// For response
	e := json.NewEncoder(w)
	var encodeError error

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		http.Error(w, gorm.ErrRecordNotFound.Error(), http.StatusNotFound)
	} else {
		encodeError = e.Encode(student)
	}

	if encodeError != nil {
		log.Printf("An error occured during encoding: %s\n", encodeError.Error())
	}
}

func newStudent(w http.ResponseWriter, r *http.Request) {
	var s models.Student
	decodeError := json.NewDecoder(r.Body).Decode(&s)
	if decodeError != nil {
		http.Error(w, decodeError.Error(), http.StatusBadRequest)
	} else {
		result := db.DB.Create(&s)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusCreated)
		}
	}
}

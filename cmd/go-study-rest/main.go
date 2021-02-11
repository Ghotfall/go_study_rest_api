package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go_study_rest_api/pkg/db"
	"go_study_rest_api/pkg/models"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	// Router setup...
	r := mux.NewRouter()
	r.HandleFunc("/students/", listStudent).Methods("GET")
	r.HandleFunc("/students/", newStudent).Methods("POST")
	r.HandleFunc("/students/{id:[0-9]+}", getStudent).Methods("GET")    // TODO: update method
	r.HandleFunc("/students/{id:[0-9]+}", updateStudent).Methods("PUT") // TODO: update method

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

	// Graceful shutdown
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
	// Find
	vars := mux.Vars(r)
	var student models.Student
	findError := models.FindFirstStudent(db.DB, &student, vars["name"])
	if findError != nil {
		http.Error(w, findError.Error(), http.StatusNotFound)
		return
	}

	// Encode
	e := json.NewEncoder(w)
	encodeError := e.Encode(student)
	if encodeError != nil {
		log.Printf("An error occured during encoding: %s\n", encodeError.Error())
		http.Error(w, encodeError.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func listStudent(w http.ResponseWriter, _ *http.Request) {
	// Get
	var students []models.Student
	getError := models.GetAllStudent(db.DB, &students)
	if getError != nil {
		http.Error(w, getError.Error(), http.StatusNotFound)
		return
	}

	// Encode
	e := json.NewEncoder(w)
	encodeError := e.Encode(students)
	if encodeError != nil {
		log.Printf("An error occured during encoding: %s\n", encodeError.Error())
		http.Error(w, encodeError.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func newStudent(w http.ResponseWriter, r *http.Request) {
	// Decode
	var s models.Student
	decodeError := json.NewDecoder(r.Body).Decode(&s)
	if decodeError != nil {
		http.Error(w, decodeError.Error(), http.StatusBadRequest)
		return
	}

	// Create
	createError := models.CreateStudent(db.DB, &s)
	if createError != nil {
		http.Error(w, createError.Error(), http.StatusInternalServerError)
		return
	}

	encodeError := json.NewEncoder(w).Encode(s)
	if encodeError != nil {
		log.Printf("An error occured during encoding: %s\n", encodeError.Error())
		http.Error(w, encodeError.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func updateStudent(w http.ResponseWriter, r *http.Request) {
	// Decode
	var s models.Student
	decodeError := json.NewDecoder(r.Body).Decode(&s)
	if decodeError != nil {
		http.Error(w, decodeError.Error(), http.StatusBadRequest)
		return
	}

	// Update
	updateError := models.SaveStudent(db.DB, &s) // TODO: update correctly
	if updateError != nil {
		http.Error(w, updateError.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

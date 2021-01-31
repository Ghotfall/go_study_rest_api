package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", mainTest)

	srv := &http.Server{
		Addr:         ":8000",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Println("Server is starting...")
	log.Fatal(srv.ListenAndServe())
}

func mainTest(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, writeErr := fmt.Fprint(w, "Hi!")
	if writeErr != nil {
		log.Printf("An error occured while executing 'mainTest' func: %s\n", writeErr.Error())
	}
}

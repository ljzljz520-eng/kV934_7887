package main

import (
	"log"
	"net/http"

	"course-trial/internal/httpapi"
	"course-trial/internal/trial"
)

func main() {
	store := trial.NewFixtureStore(trial.FixturePages())
	server := &http.Server{Addr: ":8080", Handler: httpapi.New(trial.NewService(store))}
	log.Printf("course trial server listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

package main

import (
	"log"
	"net/http"
)

type Config struct {
}

const webPort = "8084"

func main() {
	app := Config{}

	log.Println("Starting mail service on port", webPort)

	server := &http.Server{
		Addr:    ":" + webPort,
		Handler: app.routes(),
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
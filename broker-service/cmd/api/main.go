package main

import (
	"fmt"
	"log"
	"os"
	"math"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "8080"

type Config struct{
	Rabbit *amqp.Connection
}

func main() {
	// try to connect to rabitmq server
	rabbitConn, err := connect()
	if err != nil {
		log.Fatalln("Failed to connect to RabbitMQ:", err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
	fmt.Println("Connected to RabbitMQ")

	app := Config{
		Rabbit: rabbitConn,
	}

	log.Printf("Starting broker service on port %s", webPort)

	// define http server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start the server
	err = server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backoff = 1 * time.Second
	var connection *amqp.Connection

	// dont continue rabbit is ready

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
		if err != nil {
			counts++
		} else {
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backoff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backoff)
		continue
	}

	return connection, nil
}

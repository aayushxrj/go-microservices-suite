package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	ampq "github.com/rabbitmq/amqp091-go"
)

func main() {
	// try to connect to rabitmq server
	rabbitConn, err := connect()
	if err != nil {
		log.Fatalln("Failed to connect to RabbitMQ:", err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
	fmt.Println("Connected to RabbitMQ")

	// start listening for messages

	// create consumer

	// watch the queue and consume events
}

func connect() (*ampq.Connection, error) {
	var counts int64
	var backoff = 1 * time.Second
	var connection *ampq.Connection

	// dont continue rabbit is ready

	for {
		c, err := ampq.Dial("amqp://guest:guest@localhost:5672/")
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

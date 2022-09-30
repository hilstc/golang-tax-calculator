package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/hilstc/golang-tax-calculator/internal/order/infra/database"
	"github.com/hilstc/golang-tax-calculator/internal/order/usecase"
	"github.com/hilstc/golang-tax-calculator/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Uses RabbitMQ to consume the application

func main() {
	// Defines the maximum number of simultaneous workers, which will open threads to access RabbitMQ
	maxWorkers := 3
	wg := sync.WaitGroup{}

	// Opens a database connection
	db, err := sql.Open("mysql", "root:root@tcp(mysql:3306)/orders")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	repository := database.NewOrderRepository(db)
	// Creates an instance of the usecase
	uc := usecase.NewCalculateFinalPriceUseCase(repository)

	// When accessing "/total", there will be a response (what will be returned to the user: if there was an error when querying the database (error 500); the value of "total") and a request (the information received from the user)
	http.HandleFunc("/total", func(w http.ResponseWriter, r *http.Request) {
		// Creates an instance of the usecase using the repository
		uc := usecase.NewGetTotalUseCase(repository)
		// "Output" returns the result of the transaction
		output, err := uc.Execute()
		// In case of error, returns an Internal Server Error message
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// If there are no errors, converts the output to a JSON.
		// "Encode" receives a writer file
		json.NewEncoder(w).Encode(output)
	})

	// In a new thread, fires up a web server on the port below, leaving the multiplexer empty
	go http.ListenAndServe(":8181", nil)

	// Opens the channel on RabbitMQ
	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()
	// Creates a channel to establish a communication between threads
	out := make(chan amqp.Delivery)

	// Starts consuming the messages on RabbitMQ with the number of workers defined in "maxWorkers"
	go rabbitmq.Consume(ch, out)

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		i := i
		go func() {
			fmt.Println("Starting worker", i)
			defer wg.Done()
			worker(out, uc, i)
		}()
	}

	// Waits until the RabbitMQ threads finish their execution
	wg.Wait()

}

func worker(deliveryMessage <-chan amqp.Delivery, uc *usecase.CalculateFinalPriceUseCase, workerId int) {
	for msg := range deliveryMessage {
		var input usecase.OrderInputDTO

		err := json.Unmarshal(msg.Body, &input)
		if err != nil {
			fmt.Println("Error unmarshalling message:", err)
		}

		input.Tax = 10.0
		_, err = uc.Execute(input)
		if err != nil {
			fmt.Println("Error unmarshalling message:", err)
		}

		msg.Ack(false)
		fmt.Println("Worker", workerId, "processed order", input.ID)
		time.Sleep(1 * time.Second)
	}
}

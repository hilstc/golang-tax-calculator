package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/hilstc/golang-tax-calculator/internal/order/infra/database"
	"github.com/hilstc/golang-tax-calculator/internal/order/usecase"
	"github.com/hilstc/golang-tax-calculator/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(mysql:3306)/orders")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	repository := database.NewOrderRepository(db)
	uc := usecase.NewCalculateFinalPriceUseCase(repository)

	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	out := make(chan amqp.Delivery)

	go rabbitmq.Consume(ch, out)
	go worker(out, uc, 1)
	go worker(out, uc, 2)
	go worker(out, uc, 3)

	// input := usecase.OrderInputDTO{
	// 	ID:    "1234",
	// 	Price: 100,
	// 	Tax:   10,
	// }

	// output, err := uc.Execute(input)

	// if err != nil {
	// 	panic(err)
	// }

	// println(output.FinalPrice)

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
	}
}

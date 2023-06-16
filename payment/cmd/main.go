package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
)

const (
	TOPIC = "send_email"
)

type service struct {
	Broker sarama.SyncProducer
}

type Transaction struct {
	ID        string `json:"id" faker:"uuid_digit"`
	Name      string `json:"name" faker:"name"`
	Email     string `json:"email" faker:"first_name"`
	Subject   string `json:"subject" faker:"sentence"`
	Type      string `json:"type" faker:"oneof: PAY"`
	Amount    string `json:"amount" faker:"amount_with_currency"`
	CreatedAt string `json:"created_at" faker:"timestamp"`
}

func main() {
	broker, err := connectBroker()
	if err != nil {
		log.Fatalln("Failed to connect to Brokers")
	}
	app := gin.New()
	svc := service{broker}
	app.POST("/pay", svc.pay)
	app.Run("0.0.0.0:8000")
}

func connectBroker() (sarama.SyncProducer, error) {
	urls := "localhost:19092"
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	return sarama.NewSyncProducer(strings.Split(urls, ","), config)
}

func (s service) pay(ctx *gin.Context) {
	// Sent into Broker
	go func() {
		trx := Transaction{}
		err := faker.FakeData(&trx)
		if err != nil {
			log.Println("Failed generate fake data")
		}
		trx.Email = fmt.Sprintf("%s@mailhog.local", trx.Email) // ? add domain
		msgByte, err := json.Marshal(trx)
		if err != nil {
			log.Println("Failed to serialize into json")
		}
		_, _, err = s.Broker.SendMessage(&sarama.ProducerMessage{
			Topic: TOPIC,
			Value: sarama.StringEncoder(string(msgByte)),
		})
		if err != nil {
			log.Printf("Failed to push message into `%s`\n", TOPIC)
		}
	}()
	// Sent Response
	ctx.JSON(http.StatusOK, map[string]string{
		"message": "Invoice sent into your Email",
	})
}

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/Shopify/sarama"
	"gopkg.in/gomail.v2"
)

const (
	TOPIC        = "send_email"
	GROUP_ID     = "user_created_consumer"
	SENDER_EMAIL = "sender@mailhog.local"
)

type Consumer struct {
	ready      chan bool
	mailDialer *gomail.Dialer
}

type Transaction struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Subject   string `json:"subject"`
	Type      string `json:"type"`
	Amount    string `json:"amount"`
	CreatedAt string `json:"created_at"`
}

func main() {
	broker, err := connectBroker()
	if err != nil {
		log.Fatalln("Failed to connect to Brokers")
	}
	d := gomail.NewDialer("0.0.0.0", 1025, "", "")
	keepRunning := true
	ctx, cancel := context.WithCancel(context.Background())
	consumer := Consumer{
		ready:      make(chan bool),
		mailDialer: d,
	}
	consumptionIsPaused := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(ctx context.Context, worker sarama.ConsumerGroup) {
		defer wg.Done()
		for {
			if err := worker.Consume(ctx, strings.Split(TOPIC, ","), &consumer); err != nil {
				log.Fatalf("Error from consumer: %v\n", err)
			}
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}(ctx, broker)

	<-consumer.ready
	log.Println("Sarama consumer up and running!...")
	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Println("terminating: via signal")
			keepRunning = false
		case <-sigusr1:
			toggleConsumptionFlow(broker, &consumptionIsPaused)
		}
	}
	// ? Gracefully shutdown
	cancel()
	wg.Wait()
	if err = broker.Close(); err != nil {
		log.Fatalf("Error closing Consumer: %v\n", err)
	}
}

func connectBroker() (sarama.ConsumerGroup, error) {
	urls := "localhost:19092"
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	conn, err := sarama.NewConsumerGroup(strings.Split(urls, ","), GROUP_ID, config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
func toggleConsumptionFlow(client sarama.ConsumerGroup, isPaused *bool) {
	if *isPaused {
		client.ResumeAll()
		log.Println("Resuming consumption")
	} else {
		client.PauseAll()
		log.Println("Pausing consumption")
	}
	*isPaused = !*isPaused
}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.ready)
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg := <-claim.Messages():
			var trx Transaction
			if err := json.Unmarshal(msg.Value, &trx); err != nil {
				log.Println("Failed to serialize data")
				continue
			}
			content, err := parseTemplate(trx.Type, trx)
			if err != nil {
				log.Println("Failed to render email")
				log.Println(err)
				continue
			}
			m := gomail.NewMessage()
			m.SetHeader("From", SENDER_EMAIL)
			m.SetHeader("To", trx.Email)
			m.SetHeader("Subject", trx.Subject)
			m.SetBody("text/html", *content)
			if err := consumer.mailDialer.DialAndSend(m); err != nil {
				session.MarkMessage(msg, "failed")
			} else {
				session.MarkMessage(msg, "sent")
			}
		case <-session.Context().Done():
			return nil
		}
	}
}

func parseTemplate(path string, data any) (*string, error) {
	t, err := template.ParseFiles("sender/templates/out/" + path + ".html")
	if err != nil {
		return nil, err
	}
	buff := &bytes.Buffer{}
	err = t.Execute(buff, data)

	if err != nil {
		return nil, err
	}
	str := buff.String()
	return &str, nil
}

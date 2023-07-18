package main

import (
	"log"
	"logs-consumer/config"
	"logs-consumer/logs_consumer"
)

func main() {
	cred := config.LoadCredentials()
	consumer, err := logs_consumer.NewLogsConsumer(cred)
	if err != nil {
		log.Fatal(err)
	}

	err = consumer.Listen()
	if err != nil {
		log.Fatal(err)
	}
}

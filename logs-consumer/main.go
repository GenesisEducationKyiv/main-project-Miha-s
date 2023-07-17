package main

import (
	"log"
	"logs-consumer/logs_consumer"
)

func main() {
	consumer, err := logs_consumer.NewLogsConsumer("admin", "VeryStroNgAndUniqUePassWord*2s", "localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}

	err = consumer.Listen()
	if err != nil {
		log.Fatal(err)
	}
}

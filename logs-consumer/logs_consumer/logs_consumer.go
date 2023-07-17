package logs_consumer

import (
	"fmt"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
)

type LogsConsumer struct {
	conn         *amqp.Connection
	ch           *amqp.Channel
	queue        amqp.Queue
	name         string
	password     string
	urlPort      string
	exchangeName string
}

func NewLogsConsumer(name, pass, urlPort string) (*LogsConsumer, error) {
	consumer := &LogsConsumer{
		name:         name,
		password:     pass,
		urlPort:      urlPort,
		exchangeName: "logs",
	}

	err := consumer.init()
	if err != nil {
		return nil, err
	}
	return consumer, nil
}

func (consumer *LogsConsumer) Listen() error {
	defer consumer.conn.Close()
	defer consumer.ch.Close()

	msgs, err := consumer.ch.Consume(
		consumer.queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "failed to create consumer")
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
	return nil
}

func (consumer *LogsConsumer) init() error {
	var err error
	consumer.conn, err = amqp.Dial(fmt.Sprintf("amqp://%v:%v@%v", consumer.name, consumer.password, consumer.urlPort))
	if err != nil {
		return errors.Wrap(err, "failed to connect to rabbit mq")
	}

	consumer.ch, err = consumer.conn.Channel()
	if err != nil {
		return errors.Wrap(err, "failed to get channel")
	}

	err = consumer.ch.ExchangeDeclare(
		consumer.exchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "failed to create exchange")
	}

	consumer.queue, err = consumer.ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)

	if err != nil {
		return errors.Wrap(err, "failed to declare queue")
	}

	return consumer.parseArguments()
}

func (consumer *LogsConsumer) parseArguments() error {
	if len(os.Args) < 2 {
		return errors.New(fmt.Sprintf("Usage: %topic [binding_key]...", os.Args[0]))
	}
	for _, topic := range os.Args[1:] {
		log.Printf("Binding queue %s to exchange %s with routing key %s",
			consumer.queue.Name, consumer.exchangeName, topic)
		err := consumer.ch.QueueBind(
			consumer.queue.Name,
			topic,
			"logs",
			false,
			nil)
		if err != nil {
			return errors.Wrap(err, "failed to bind queue")
		}
	}

	return nil
}

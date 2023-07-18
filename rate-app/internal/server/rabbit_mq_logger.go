package server

import (
	"btc-test-task/internal/common/configuration/config"
	"btc-test-task/internal/common/configuration/logger"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQLogger struct {
	conn         *amqp.Connection
	ch           *amqp.Channel
	exchangeName string
	errorsTopic  string
	successTopic string
}

type RequestInfo struct {
	Method   string
	Endpoint string
	Duration time.Duration
	Status   int
}

func NewRabbitMQLogger(conf *config.Config) (*RabbitMQLogger, error) {
	log := &RabbitMQLogger{}
	log.exchangeName = "logs"
	log.errorsTopic = "error"
	log.successTopic = "ok"
	err := log.init(conf)
	if err != nil {
		return nil, err
	}
	return log, nil
}

type MiddlewareFunc func(http.Handler) http.Handler

func (log *RabbitMQLogger) MiddlewareLogger() MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			handler.ServeHTTP(ww, r)
			info := RequestInfo{
				Method:   r.Method,
				Endpoint: r.RequestURI,
				Duration: time.Since(start),
				Status:   ww.Status(),
			}
			err := log.log(ww.Status(), info.Format())
			if err != nil {
				logger.Log.Warnf("Failed to send log %v", err)
			}
		}
		return http.HandlerFunc(fn)
	}
}

func (log *RabbitMQLogger) Close() {
	defer log.conn.Close()
	defer log.ch.Close()
}

func (req *RequestInfo) Format() string {
	str := fmt.Sprintf("%v %v - %v %v", req.Method, req.Endpoint, req.Status, req.Duration)
	return str
}

func (log *RabbitMQLogger) init(conf *config.Config) error {
	var err error
	log.conn, err = amqp.Dial(fmt.Sprintf("amqp://%v:%v@%v",
		conf.RabbitUsername, conf.RabbitPassword, conf.RabbitUrlPort))
	if err != nil {
		return err
	}

	log.ch, err = log.conn.Channel()
	if err != nil {
		return err
	}

	err = log.ch.ExchangeDeclare(
		log.exchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	return nil
}

func (log *RabbitMQLogger) log(status int, str string) error {
	var topic string
	if status >= 400 {
		topic = log.errorsTopic
	} else {
		topic = log.successTopic
	}

	return log.send(topic, str)
}

func (log *RabbitMQLogger) send(topic, str string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := log.ch.PublishWithContext(ctx,
		log.exchangeName,
		topic,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(str),
		})

	if err != nil {
		return err
	}
	return nil
}

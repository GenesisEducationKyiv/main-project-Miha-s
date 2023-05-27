package main

import (
	"btc-test-task/internal/config"
	"btc-test-task/internal/email_sender"
	"btc-test-task/internal/emails_storage"
	"btc-test-task/internal/logger"
	"btc-test-task/internal/server"
	"os"
	"os/signal"
	"syscall"
	// "crypto/tls"
	// gomail "gopkg.in/gomail.v2"
)

func main() {
	done := make(chan error, 1)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	var conf config.Config
	conf.EmailServiceUrl = "smtp.gmail.com"
	conf.EmailServicePort = 587
	conf.EmailToSendFrom = "genesis.task.mykhailo@gmail.com"
	conf.EmailToSendFromPassword = "lgjvxuatvpeislnb"
	conf.EmailSubject = "Would you like to get a new course of BTC?"
	conf.EmailStoragePath = "/home/mstatnik/golang/btc-test-task"
	email_sender.Init(&conf)
	logger.Init(&conf)
	emails_storage.Init(&conf)
	defer emails_storage.Close()
	var server server.Server
	server.Init(&conf)

	go server.Run()

	select {
	case <-signals:
	case <-done:
	}
}

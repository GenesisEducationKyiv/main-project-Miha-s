package main

import (
	"btc-test-task/internal/config"
	"btc-test-task/internal/lifecycle"
)

func main() {

	var conf config.Config
	conf.EmailServiceUrl = "smtp.gmail.com"
	conf.EmailServicePort = 587
	conf.EmailToSendFrom = "genesis.task.mykhailo@gmail.com"
	conf.EmailToSendFromPassword = "lgjvxuatvpeislnb"
	conf.EmailSubject = "Would you like to get a new course of BTC?"
	conf.EmailStoragePath = "/home/mstatnik/golang/btc-test-task"
	conf.CoinAPIUrl = "https://rest.coinapi.io/v1/exchangerate/BTC/UAH"

	var lifecycle lifecycle.Lifecycle
	lifecycle.Init(&conf)
	lifecycle.Run()

}

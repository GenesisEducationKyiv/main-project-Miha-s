package main

import (
	"btc-test-task/internal/config"
	"btc-test-task/internal/lifecycle"
	"btc-test-task/internal/logger"
)

func main() {
	var conf config.Config
	err := conf.LoadFromENV()
	if err != nil {
		panic(1)
	}

	var lifecycle lifecycle.Lifecycle
	err = lifecycle.Init(&conf)
	if err != nil {
		logger.LogError(err)
		return
	}

	err = lifecycle.Run()
	if err != nil {
		logger.LogError(err)
		return
	}
}

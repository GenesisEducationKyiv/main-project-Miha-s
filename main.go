package main

import (
	"btc-test-task/internal/configuration/config"
	"btc-test-task/internal/configuration/logger"
	"btc-test-task/internal/lifecycle"
)

func main() {
	var conf config.Config
	err := conf.LoadFromENV(".env")
	if err != nil {
		panic(1)
	}

	var lifeCycle lifecycle.Lifecycle
	err = lifeCycle.Init(&conf)
	if err != nil {
		logger.Log.Error(err)
		return
	}

	err = lifeCycle.Run()
	if err != nil {
		logger.Log.Error(err)
		return
	}
}

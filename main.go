package main

import (
	"btc-test-task/internal/config"
	"btc-test-task/internal/lifecycle"
)

func main() {

	var conf config.Config
	conf.LoadFromENV()

	var lifecycle lifecycle.Lifecycle
	lifecycle.Init(&conf)
	lifecycle.Run()

}

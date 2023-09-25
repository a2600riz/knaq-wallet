package main

import (
	"knaq-wallet/config"
	"knaq-wallet/controller"
	"knaq-wallet/database"
)

func main() {
	config.New()
	database.New()
	controller.Start()
}

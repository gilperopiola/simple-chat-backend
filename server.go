package main

import (
	"log"
	"os"
)

var config MyConfig
var database MyDatabase
var router MyRouter
var bot MyBot
var channel chan float64

func main() {
	config.Setup()

	database.Setup()
	database.Seed()
	defer database.Close()

	router.Setup()

	bot.Setup()
	defer bot.Conn.Close()
	defer bot.Channel.Close()
	go bot.Run()

	log.Println("server started")
	router.Run(":" + os.Getenv("PORT"))
}

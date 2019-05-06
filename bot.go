package main

import (
	"encoding/csv"
	"log"
	"net/http"
	"strconv"

	"github.com/streadway/amqp"
)

type BotActions interface {
	Setup()
	Run()
}

type MyBot struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Queue   amqp.Queue
}

func (bot *MyBot) Setup() {
	channel = make(chan float64)

	var err error
	bot.Conn, err = amqp.Dial("amqp://" + config.RABBITMQ.USERNAME + ":" + config.RABBITMQ.PASSWORD + "@" + config.RABBITMQ.HOST + ":" + config.RABBITMQ.PORT + "/")
	if err != nil {
		log.Fatalf("error connecting to RabbitMQ server: %v", err)
	}

	bot.Channel, err = bot.Conn.Channel()
	if err != nil {
		log.Fatalf("error connecting to RabbitMQ channel: %v", err)
	}

	bot.Queue, err = bot.Channel.QueueDeclare("hello", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("error connecting to RabbitMQ queue: %v", err)
	}
}

func (bot *MyBot) Run() {
	msgs, err := bot.Channel.Consume(bot.Queue.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Printf("error consuming RabbitMQ channel: %v", err)
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			resp, err := http.Get("https://stooq.com/q/l/?s=" + string(d.Body) + ".us&f=sd2t2ohlcv&h&e=csv")
			if err != nil {
				log.Printf("error accessing stooq API: %v", err)
			}
			defer resp.Body.Close()

			reader := csv.NewReader(resp.Body)
			reader.Comma = ','
			data, err := reader.ReadAll()
			if err != nil {
				log.Printf("error reading csv: %v", err)
			}

			stockPrice, _ := strconv.ParseFloat(data[1][6], 64)
			channel <- stockPrice
		}
	}()
	<-forever
}

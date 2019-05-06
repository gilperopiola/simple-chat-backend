package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

type Command struct {
	Command string `json:"command"`
}

func ExecuteCommand(c *gin.Context) {
	var command Command
	c.BindJSON(&command)
	command.Command = strings.ToUpper(command.Command)

	if !strings.Contains(command.Command, "STOCK") {
		c.JSON(http.StatusBadRequest, "command not recognized")
		return
	}

	split := strings.SplitAfter(command.Command, "=")
	if len(split) != 2 {
		c.JSON(http.StatusBadRequest, "invalid format")
		return
	}

	stockPrice, err := sendRabbitMQMessage(split[1])
	if err != nil {
		c.JSON(http.StatusInternalServerError, "rabbitMQ error")
		return
	}

	if stockPrice == 0 {
		c.JSON(http.StatusBadRequest, "no stock data for company "+split[1])
		return
	}

	c.JSON(http.StatusOK, strings.ToUpper(split[1])+" quote is $"+fmt.Sprintf("%.2f", stockPrice)+" per share")
}

func sendRabbitMQMessage(msg string) (float64, error) {
	conn, err := amqp.Dial("amqp://" + config.RABBITMQ.USERNAME + ":" + config.RABBITMQ.PASSWORD + "@" + config.RABBITMQ.HOST + ":" + config.RABBITMQ.PORT + "/")
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return 0, err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("hello", false, false, false, false, nil)
	if err != nil {
		return 0, err
	}

	err = ch.Publish("", q.Name, false, false, amqp.Publishing{ContentType: "text/plain", Body: []byte(msg)})
	if err != nil {
		return 0, err
	}

	stockPrice := <-channel

	return stockPrice, nil
}

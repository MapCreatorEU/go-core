package rabbitmq

import (
	"github.com/streadway/amqp"
	"fmt"
	"encoding/json"
	"github.com/MapOnline/go-core/logger"
)

type RabbitMQConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Server string `json:"server"`
	Port int `json:"port"`
	Host string `json:"host"`
	Queue string `json:"queue"`
}

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel *amqp.Channel
	Queue amqp.Queue
}

func (MessageQueue *RabbitMQ) ConnectWithQueue(Config RabbitMQConfig, Queue string) {
	var err error

	MessageQueue.Connection, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d%s",
		Config.Username,
		Config.Password,
		Config.Server,
		Config.Port,
		Config.Host,
	))
	logger.FailOnError(err, "Failed to connect to server")

	MessageQueue.Channel, err = MessageQueue.Connection.Channel()
	logger.FailOnError(err, "Failed to open channel")

	MessageQueue.Queue, err =  MessageQueue.Channel.QueueDeclare(
		Queue,
		true,
		false,
		false,
		false,
		nil,
	)
	logger.FailOnError(err, "Failed to declare queue")
}

func (MessageQueue *RabbitMQ) Connect(Config RabbitMQConfig) {
	MessageQueue.ConnectWithQueue(Config, Config.Queue)
}

func (MessageQueue *RabbitMQ) Disconnect() {
	MessageQueue.Channel.Close()
	MessageQueue.Connection.Close()
}

func (MessageQueue *RabbitMQ) SendJSON(Model interface{}) {
	var err error
	var body []byte

	body, err = json.Marshal(Model)

	if len(body) == 0 || body == nil{
		logger.Log.Println("Cannot parse model to JSON byte array empty")
		return
	}

	if err != nil {
		logger.LogOnError(err, "Cannot parse model to JSON")
		return
	}

	err = MessageQueue.Channel.Publish(
		"",
		MessageQueue.Queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: body,
		})

	logger.LogOnError(err, "Cannot send model to server")
}

func (MessageQueue *RabbitMQ) Send(Message string) {
	var err error
	err = MessageQueue.Channel.Publish(
		"",
		MessageQueue.Queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: []byte(Message),
		},
	)
	logger.LogOnError(err, "Cannot send message to server")
}
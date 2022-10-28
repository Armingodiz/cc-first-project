package brokerService

import (
	"encoding/json"
	"cc-first-project/user-service/models"

	"github.com/streadway/amqp"
)

type BrokerService interface {
	Publish(ad models.Advertisement) error
}

func NewBrokerService() BrokerService {
	return &RabbitMQBrokerService{}
}

type RabbitMQBrokerService struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func (r *RabbitMQBrokerService) Publish(ad models.Advertisement) error {
	// Define RabbitMQ server URL.
	amqpServerURL := "amqps://dfyoughc:4RKqichqlzRebd-zntErvAecOVCM5rkM@kangaroo.rmq.cloudamqp.com/dfyoughc"
	amqpQueueName := "advertisements"
	// Create a new RabbitMQ connection.
	connectRabbitMQ, err := amqp.Dial(amqpServerURL)
	if err != nil {
		return err
	}
	defer connectRabbitMQ.Close()

	// Let's start by opening a channel to our RabbitMQ
	// instance over the connection we have already
	// established.
	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		return err
	}
	defer channelRabbitMQ.Close()
	// With the instance and declare Queues that we can
	// publish and subscribe to.
	_, err = channelRabbitMQ.QueueDeclare(
		amqpQueueName, // queue name
		true,          // durable ==> if server restarts, messages will be there
		false,         // auto delete
		false,         // exclusive
		false,         // no wait
		nil,           // arguments
	)
	if err != nil {
		return err
	}
	adBytes, err := json.Marshal(ad)
	if err != nil {
		return err
	}
	message := amqp.Publishing{
		DeliveryMode: amqp.Persistent, // Persistent ==> if server restarts, messages will be there
		ContentType:  "application/json",
		Body:         adBytes,
	}

	// Attempt to publish a message to the queue.
	return channelRabbitMQ.Publish(
		"",            // exchange
		amqpQueueName, // queue name
		false,         // mandatory
		false,         // immediate
		message,       // message to publish
	)
}

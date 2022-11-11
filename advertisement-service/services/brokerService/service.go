package brokerService

import (
	"encoding/json"
	"log"
	"cc-first-project/advertisement-service/models"

	"github.com/streadway/amqp"
)

type BrokerService interface {
	StartConsuming() (chan models.Advertisement, chan error, error)
	CloseConnection() error
	CloseChannel() error
}

func NewBrokerService() BrokerService {
	return &RabbitMQBrokerService{}
}

type RabbitMQBrokerService struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func (r *RabbitMQBrokerService) StartConsuming() (advertisementChann chan models.Advertisement, errChann chan error, err error) {
	amqpServerURL := "amqps://dfyoughc:4RKqichqlzRebd-zntErvAecOVCM5rkM@kangaroo.rmq.cloudamqp.com/dfyoughc"
	amqpQueueName := "advertisements"
	connectRabbitMQ, err := amqp.Dial(amqpServerURL)
	if err != nil {
		return
	}
	r.Connection = connectRabbitMQ
	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		return
	}
	r.Channel = channelRabbitMQ
	queue, err := r.Channel.QueueDeclare(
		amqpQueueName, // name
		true,          // durable (will survive server restarts)
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		return
	}
	err = r.Channel.Qos( // it is for fair dispatch and means if there is no free workers, the message will be put in the queue and will be delivered to the next worker.
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return
	}

	// Subscribing to QueueService1 for getting messages.
	messages, err := channelRabbitMQ.Consume(
		queue.Name, // queue name
		"",         // consumer
		false,      // auto-ack == > message.Ack(false) ==> message will be removed
		false,      // exclusive
		false,      // no local
		false,      // no wait
		nil,        // arguments
	)
	if err != nil {
		return
	}

	advertisementChann = make(chan models.Advertisement, 10)
	errChann = make(chan error, 2)
	go func() {
		for message := range messages {
			var advertisement models.Advertisement
			err := json.Unmarshal(message.Body, &advertisement)
			if err != nil {
				log.Println("Error:", err)
				errChann <- err
				return
			}
			message.Ack(false)
			advertisementChann <- advertisement
		}
	}()
	return
}

func (r *RabbitMQBrokerService) CloseConnection() error {
	return r.Connection.Close()
}

func (r *RabbitMQBrokerService) CloseChannel() error {
	return r.Channel.Close()
}

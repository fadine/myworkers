package queue

import (
	"encoding/json"

	"github.com/fadine/myworkers/global"
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

type RabbitMessage struct {
	amqp.Delivery
}

func (r RabbitMQ) Close() {

	if r.ch != nil {
		r.ch.Close()
	}

	if r.conn != nil {
		r.conn.Close()
	}
}

func (r RabbitMQ) Send(queueName string, body string) {
	ch := r.getChannel()
	q, _ := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	_ = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType:     "application/json",
			ContentEncoding: "",
			Body:            []byte(body),
		})

}

func (r RabbitMQ) GetMessageFromQueue(queueName string) <-chan IQueueMessage {

	ch := r.getChannel()
	q, _ := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	_ = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	msgs, _ := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	delivery := make(chan IQueueMessage)
	go func() {

		for {
			select {
			case <-global.GracefulStop:
				r.Close()
				close(delivery)
				return
			case msg := <-msgs:
				message := RabbitMessage{Delivery: msg}
				if message.IsValid() {
					delivery <- message
				} else {
					msg.Ack(false)
				}
			}
		}
	}()

	return (<-chan IQueueMessage)(delivery)
}

func (r *RabbitMQ) getConnect() *amqp.Connection {

	if r.conn == nil {

		//var err error
		host, _ := global.Cfg.String("rabbitmq")

		r.conn, _ = amqp.Dial(host)
	}
	return r.conn
}

func (r *RabbitMQ) getChannel() *amqp.Channel {

	conn := r.getConnect()
	if r.ch == nil {

		ch, _ := conn.Channel()
		r.ch = ch
	}

	return r.ch
}

func (m RabbitMessage) GetAction() string {

	var payload MessagePayload
	err := json.Unmarshal(m.GetBody(), &payload)
	if err != nil {
		panic(err)
	}

	return payload.Action
}

func (m RabbitMessage) GetBody() []byte {
	return m.Body
}

func (m RabbitMessage) IsValid() bool {
	return m.ContentType == "application/json"
}

func (m RabbitMessage) GetId() string {
	return m.MessageId
}

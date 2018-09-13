package queue

import (
	"fmt"
)

type MessagePayload struct {
	Action string `json:"action"`
}

var drivers = map[string]IQueue{
	"rabbitmq": RabbitMQ{},
}

func GetService() IQueue {

	driver := "rabbitmq"
	service, ok := drivers[driver]
	if !ok {
		panic(fmt.Sprintf("Queue driver %s is not found", driver))
	}

	return service
}

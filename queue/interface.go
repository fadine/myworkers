package queue

type IQueue interface {
	GetMessageFromQueue(string) <-chan IQueueMessage
	Close()
	Send(string, string)
}

type IQueueMessage interface {
	Ack(multiple bool) error
	Reject(requeue bool) error
	Nack(multiple, requeue bool) error
	GetBody() []byte
	GetAction() string
	GetId() string
}

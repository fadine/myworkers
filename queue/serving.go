package queue

type ServingMessage struct {
	Action string
}

func (s ServingMessage) Ack(multiple bool) error {
	return nil
}

func (s ServingMessage) Reject(requeue bool) error {
	return nil
}

func (s ServingMessage) Nack(multiple, requeue bool) error {
	return nil
}

func (s ServingMessage) GetBody() []byte {
	return nil
}

func (s ServingMessage) GetAction() string {
	return s.Action
}

func (s ServingMessage) GetId() string {
	return ""
}

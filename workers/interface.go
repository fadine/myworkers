package workers

import (
	"sync"

	"github.com/fadine/myworkers/queue"
)

type IHandler interface {
	Process(msg queue.IQueueMessage, group *sync.WaitGroup)
}

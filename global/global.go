package global

import (
	"github.com/olebedev/config"
	"github.com/satori/go.uuid"
)

var (
	GracefulStop = make(chan bool)
	MustStop     = make(chan bool)
	SafelyQuit   = false

	WorkerId, _ = uuid.NewV4()

	Cfg config.Config
)

package workers

import (
	"fmt"
	"reflect"

	"github.com/fadine/myworkers/workers/handlers"
)

var mapHandlers = map[string]reflect.Type{
	"EmailEngine": reflect.TypeOf(handlers.EmailEngine{}),
}

func GetHandler(action string) IHandler {
	fmt.Println("action ==> ", action)
	handler, ok := mapHandlers[action]
	if !ok {
		panic(fmt.Sprintf("Handler %s is not valid", action))
	}

	return reflect.New(handler).Elem().Interface().(IHandler)
}

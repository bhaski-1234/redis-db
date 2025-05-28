package dispatcher

import (
	"errors"
	"strings"

	"github.com/bhaski-1234/redis-db/internal/command"
)

type HandlerFunc func(args []string) (interface{}, error)

type Dispatcher struct {
	handlers map[string]HandlerFunc
}

func NewDispatcher() *Dispatcher {
	d := &Dispatcher{
		handlers: make(map[string]HandlerFunc),
	}

	// Register commands
	d.Register("PING", command.HandlePing)
	d.Register("GET", command.HandleGet)
	d.Register("SET", command.HandleSet)

	return d
}

func (d *Dispatcher) Register(cmd string, handler HandlerFunc) {
	d.handlers[strings.ToUpper(cmd)] = handler
}

func (d *Dispatcher) Execute(cmd string, args []string) (interface{}, error) {
	handler, exists := d.handlers[strings.ToUpper(cmd)]
	if !exists {
		return nil, errors.New("ERR unknown command '" + cmd + "'")
	}

	return handler(args)
}

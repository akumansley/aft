package bus

import (
	"context"
	"reflect"
	"sync"
)

// https://play.golang.org/p/uversZMgwT
type EventBus struct {
	handlers map[reflect.Type][]reflect.Value
	lock     sync.RWMutex
}

func New() *EventBus {
	return &EventBus{
		make(map[reflect.Type][]reflect.Value),
		sync.RWMutex{},
	}
}

func (bus *EventBus) addHandler(t reflect.Type, fn reflect.Value) {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	handlers, ok := bus.handlers[t]
	if !ok {
		handlers = make([]reflect.Value, 0)
	}
	bus.handlers[t] = append(handlers, fn)
}

func (bus *EventBus) RegisterHandlers(fns []interface{}) {
	for _, f := range fns {
		bus.RegisterHandler(f)
	}
}

func (bus *EventBus) RegisterHandler(fn interface{}) {
	v := reflect.ValueOf(fn)
	def := v.Type()

	// the message handler must have a single parameter
	if def.NumIn() != 1 {
		panic("Handler must have a single argument")
	}
	// find out the handler argument type
	argument := def.In(0)

	bus.addHandler(argument, v)
}

func (bus *EventBus) Publish(ev interface{}) error {
	bus.lock.RLock()
	defer bus.lock.RUnlock()

	t := reflect.TypeOf(ev)

	handlers, ok := bus.handlers[t]
	if !ok {
		return nil
	}

	args := [...]reflect.Value{reflect.ValueOf(ev)}
	for _, fn := range handlers {
		fn.Call(args[:])
	}
	return nil
}

// Help the server stash the bus in a context
type key int

const busKey key = 0

func WithBus(ctx context.Context, bus *EventBus) context.Context {
	return context.WithValue(ctx, busKey, bus)
}
func FromContext(ctx context.Context) (bus *EventBus, ok bool) {
	bus, ok = ctx.Value(busKey).(*EventBus)
	return
}

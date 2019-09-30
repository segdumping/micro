package configcenter

import (
	"reflect"
	"sync"
)

type Watch interface {
	Watch(map[string][]byte)
}

type dispatch struct {
	sync.RWMutex
	listeners map[string][]Watch
}

func newDispatch() *dispatch {
	return &dispatch{
		listeners: map[string][]Watch{},
	}
}

func (d *dispatch) AddListen(event string, listener Watch) {
	d.Lock()
	d.listeners[event] = append(d.listeners[event], listener)
	d.Unlock()
}

func (d *dispatch) RemoveListen(event string, listener Watch) {
	d.Lock()
	defer d.Unlock()
	listeners, ok := d.listeners[event]
	if !ok {
		return
	}

	for i, v := range listeners {
		if !reflect.DeepEqual(listener, v) {
			continue
		}

		listeners = append(listeners[:i], listeners[i+1:]...)
		d.listeners[event] = listeners
		break
	}
}

func (d *dispatch) fire(e *event) {
	d.RLock()
	defer d.RUnlock()
	listeners, ok := d.listeners[e.key]
	if !ok {
		return
	}

	for _, v := range listeners {
		v.Watch(map[string][]byte{e.key: e.value})
	}
}

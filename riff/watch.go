package riff

import (
	"sync"
)

type WatchType int

const (
	NodeChanged WatchType = 1 << iota
	ServiceChanged
)

type WatchParam struct {
	Name string
	WatchType
}

type WatchHandler interface {
	GetParam() *WatchParam
	HandleWatch()
}

type Watch struct {
	sync.Mutex
	handlers map[WatchHandler]struct{}
}

func NewWatch() *Watch {
	return &Watch{
		handlers: make(map[WatchHandler]struct{}),
	}
}

func (w *Watch) RegisterHandler(wh WatchHandler) {
	w.Lock()
	defer w.Unlock()

	// Do nothing if already registered
	if _, ok := w.handlers[wh]; ok {
		return
	}

	// Register
	w.handlers[wh] = struct{}{}

	// Dispatch once
	wh.HandleWatch()
}

func (w *Watch) DeregisterHandler(wh WatchHandler) {
	w.Lock()
	defer w.Unlock()
	delete(w.handlers, wh)
}

func (w *Watch) Dispatch(param WatchParam) {
	for wh := range w.handlers {
		whParam := wh.GetParam()
		if whParam.Name == param.Name && whParam.WatchType == param.WatchType {
			wh.HandleWatch()
		}
	}
}

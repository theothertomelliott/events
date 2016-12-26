package events

import (
	"errors"
	"sync"
)

// FanOut forwards everything sent to an input channel to a set of listener channels
type FanOut struct {
	lag         int
	outChannels map[<-chan interface{}]chan interface{}
	mutex       *sync.Mutex
	closed      bool
}

// NewFanOut creates an instance of FanOut which accepts input on the provided channel,
// with listener channels with a buffer size defined by lag.
func NewFanOut(inChannel <-chan interface{}, lag int) *FanOut {
	fanOut := &FanOut{
		lag:         lag,
		outChannels: make(map[<-chan interface{}]chan interface{}),
		mutex:       &sync.Mutex{},
	}
	go func() {
		for value := range inChannel {
			fanOut.mutex.Lock()
			for _, c := range fanOut.outChannels {
				c <- value
			}
			fanOut.mutex.Unlock()
		}
		fanOut.mutex.Lock()
		fanOut.closed = true
		for _, c := range fanOut.outChannels {
			close(c)
		}
		fanOut.mutex.Unlock()
	}()
	return fanOut
}

// Listen creates a channel to listen for output from this instance of FanOut
func (f *FanOut) Listen() (<-chan interface{}, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if f.closed {
		return nil, errors.New("input channel already closed")
	}

	newChan := make(chan interface{}, f.lag)
	f.outChannels[newChan] = newChan
	return newChan, nil
}

// StopListening closes a listener channel previously created by Listen. If this channel
// has already been closed or was not created by Listen, an error will be returned.
func (f *FanOut) StopListening(c <-chan interface{}) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if f.closed {
		return errors.New("input channel already closed")
	}

	if _, ok := f.outChannels[c]; ok {
		delete(f.outChannels, c)
		return nil
	}

	return errors.New("channel not found")
}

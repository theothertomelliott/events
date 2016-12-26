package events

import (
	"errors"
	"sync"
)

type FanOut struct {
	lag         int
	outChannels []chan interface{}
	mutex       *sync.Mutex
	closed      bool
}

func NewFanOut(inChannel <-chan interface{}, lag int) *FanOut {
	fanOut := &FanOut{
		lag:   lag,
		mutex: &sync.Mutex{},
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

func (f *FanOut) Listen() (<-chan interface{}, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if f.closed {
		return nil, errors.New("input channel already closed")
	}

	newChan := make(chan interface{}, f.lag)
	f.outChannels = append(f.outChannels, newChan)

	return newChan, nil
}

func (f *FanOut) StopListening(c <-chan interface{}) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if f.closed {
		return errors.New("input channel already closed")
	}

	for i, co := range f.outChannels {
		if co == c {
			if i < len(f.outChannels)-1 {
				f.outChannels[i] = f.outChannels[len(f.outChannels)-1]
			}
			f.outChannels = f.outChannels[:len(f.outChannels)-1]
			close(co)
			return nil
		}
	}

	return errors.New("channel not found")
}

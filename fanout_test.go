package events

import "testing"

func TestCreateFanout(t *testing.T) {
	inChannel := make(chan interface{})
	f := NewFanOut(inChannel, 1)
	c, err := f.Listen()
	if err != nil {
		t.Error(err)
	}
	inChannel <- "test"

	close(inChannel)

	out, ok := <-c
	if !ok {
		t.Error("Nothing on channel")
		return
	}
	if "test" != out {
		t.Errorf("Expected 'test', got '%v'", out)
	}
}

func TestStopListen(t *testing.T) {
	inChannel := make(chan interface{})
	f := NewFanOut(inChannel, 1)
	c, err := f.Listen()
	if err != nil {
		t.Error(err)
	}
	err = f.StopListening(c)
	if err != nil {
		t.Error(err)
	}
	err = f.StopListening(c)
	if err == nil {
		t.Error("Expected an error")
	}
	if err.Error() != "channel not found" {
		t.Error("Unexpected error:", err)
	}
	close(inChannel)
}

func TestStopListenMultiple(t *testing.T) {
	inChannel := make(chan interface{})
	f := NewFanOut(inChannel, 1)
	c, err := f.Listen()
	if err != nil {
		t.Error(err)
	}
	_, err = f.Listen()
	if err != nil {
		t.Error(err)
	}
	err = f.StopListening(c)
	if err != nil {
		t.Error(err)
	}
	err = f.StopListening(c)
	if err == nil {
		t.Error("Expected an error")
	}
	if err.Error() != "channel not found" {
		t.Error("Unexpected error:", err)
	}
	close(inChannel)
}

func TestCloseChannel(t *testing.T) {

	inChannel := make(chan interface{})
	f := NewFanOut(inChannel, 1)
	c, err := f.Listen()
	if err != nil {
		t.Error(err)
	}

	close(inChannel)

	_ = <-c

	_, err = f.Listen()
	if err == nil {
		t.Error("Expected an error")
	} else if err.Error() != "input channel already closed" {
		t.Error("Unexpected error:", err)
	}

	err = f.StopListening(c)
	if err == nil {
		t.Error("Expected an error")
	} else if err.Error() != "input channel already closed" {
		t.Error("Unexpected error:", err)
	}
}

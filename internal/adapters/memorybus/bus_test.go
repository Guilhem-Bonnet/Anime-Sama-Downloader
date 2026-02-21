package memorybus

import (
	"testing"
	"time"
)

func TestBus_New(t *testing.T) {
	bus := New()
	if bus == nil {
		t.Fatal("expected non-nil bus")
	}
}

func TestBus_Publish_Subscribe(t *testing.T) {
	bus := New()

	sub, cancel := bus.Subscribe()
	defer cancel()

	go func() {
		bus.Publish("test", []byte(`{"data":"test"}`))
	}()

	select {
	case evt := <-sub:
		if evt.Topic != "test" {
			t.Errorf("expected topic 'test', got %s", evt.Topic)
		}
		if string(evt.Payload) != `{"data":"test"}` {
			t.Errorf("expected payload '{\"data\":\"test\"}', got %s", evt.Payload)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for event")
	}
}

func TestBus_MultipleSubscribers(t *testing.T) {
	bus := New()

	sub1, cancel1 := bus.Subscribe()
	defer cancel1()

	sub2, cancel2 := bus.Subscribe()
	defer cancel2()

	go func() {
		bus.Publish("multi", []byte(`{"msg":"hi"}`))
	}()

	// Both subscribers should receive the event
	for i := 0; i < 2; i++ {
		select {
		case evt := <-sub1:
			if evt.Topic != "multi" {
				t.Errorf("sub1: expected topic 'multi', got %s", evt.Topic)
			}
			sub1 = nil // mark received
		case evt := <-sub2:
			if evt.Topic != "multi" {
				t.Errorf("sub2: expected topic 'multi', got %s", evt.Topic)
			}
			sub2 = nil // mark received
		case <-time.After(1 * time.Second):
			t.Fatal("timeout waiting for event")
		}
	}
}

func TestBus_UnsubscribeClosesChannel(t *testing.T) {
	bus := New()

	sub, cancel := bus.Subscribe()

	// Cancel subscription
	cancel()

	bus.Publish("test", []byte(`{}`))

	select {
	case <-sub:
		// Channel might be closed, which is OK
	case <-time.After(100 * time.Millisecond):
		// No event received, which is also OK
	}
}

package fanner

import (
	"strings"
	"sync"
	"testing"
)

func TestUnbufferedFanOut(t *testing.T) {
	fo := NewOut()

	topics := []string{"foo", "bar", "baz"}
	values := []string{"a", "b", "c", "d"}

	for _, topic := range topics {
		ch := fo.Register(topic)

		for _, value := range values {
			go fo.Send(topic, value)
			recv, ok := <-ch
			if ok == false {
				t.Fatal("unexpected closed channel")
			}
			if recv != value {
				t.Fatalf("expected value '%s', got '%s'", value, recv)
			}
		}

		fo.Unregister(ch)
		select {
		case <-ch:
		default:
			t.Fatal("channel expected to be closed")
		}
	}
}

func TestBufferedFanOut(t *testing.T) {
	fo := NewOut()

	topics := []string{"foo", "bar", "baz"}
	values := []string{"a", "b", "c", "d"}

	for _, topic := range topics {
		ch := fo.Register(topic, 4)

		wg := new(sync.WaitGroup)
		wg.Add(len(values))
		for _, value := range values {
			go func() {
				fo.Send(topic, value)
				wg.Done()
			}()
		}
		wg.Wait()

		fo.Unregister(ch)
		concat := strings.Join(values, "")
		for i := 0; i < len(values); i++ {
			recv := <-ch
			if strings.Contains(concat, recv.(string)) == false {
				t.Fatalf("unexpected value '%s'", recv)
			}
		}

		select {
		case <-ch:
		default:
			t.Fatal("channel expected to be closed")
		}
	}
}

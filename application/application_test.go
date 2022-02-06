package application

import (
	"sync"
	"testing"
	"time"
)

func TestDefer(t *testing.T) {
	app := New()

	var deferred bool
	app.Defer(func() { deferred = true })

	go func() {
		time.Sleep(time.Second)
		app.Stop()
	}()

	app.Run()

	if deferred == false {
		t.Error("Failed to invoke deferred functions.")
	}
}

func TestStop(t *testing.T) {
	app := New()

	var deferred bool
	app.Defer(func() { deferred = true })

	var stop bool
	app.OnStop(func() { stop = true })

	go func() {
		time.Sleep(time.Second)
		app.Stop()
	}()

	app.Run()

	if deferred == false {
		t.Error("Failed to invoke deferred functions.")
	}

	if stop == false {
		t.Error("Failed to invoke stop functions.")
	}
}

func TestForcedStop(t *testing.T) {
	app := New()

	var deferred bool
	app.Defer(func() { deferred = true })

	wg := new(sync.WaitGroup)
	wg.Add(1)
	app.OnStop(func() { wg.Wait() }) // this will block indefinitely

	var stop bool
	app.OnStop(func() { stop = true })

	var forcedStop bool
	app.OnForcedStop(func() { forcedStop = true })

	go func() {
		time.Sleep(time.Second)
		app.Stop()
		app.Stop()
	}()

	app.Run()

	if deferred == false {
		t.Error("Failed to invoke deferred functions.")
	}

	if stop {
		t.Error("Unexpected invocation of stop functions.")
	}

	if forcedStop == false {
		t.Error("Failed to invoke forced stop functions.")
	}
}

func TestBlockedForcedStop(t *testing.T) {
	app := New()

	var deferred bool
	app.Defer(func() { deferred = true })

	wg := new(sync.WaitGroup)
	wg.Add(1)
	app.OnStop(func() { wg.Wait() }) // this will block indefinitely

	var stop bool
	app.OnStop(func() { stop = true })

	wg.Add(1)
	app.OnForcedStop(func() { wg.Wait() }) // this will block indefinitely

	var forcedStop bool
	app.OnForcedStop(func() { forcedStop = true })

	go func() {
		time.Sleep(time.Second)
		app.Stop()
		app.Stop()
		app.Stop()
	}()

	app.Run()

	if deferred == false {
		t.Error("Failed to invoke deferred functions.")
	}

	if stop {
		t.Error("Unexpected invocation of stop functions.")
	}

	if forcedStop {
		t.Error("Unexpected invocation of forced stop functions.")
	}
}

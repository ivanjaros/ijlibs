//go:build !windows
// +build !windows

package application

import "testing"

func TestSignalStop(t *testing.T) {
	signals := []os.Signal{
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGQUIT,
		syscall.SIGKILL,
	}

	for _, sig := range signals {
		t.Log("Testing ", sig, "signal.")
		app := New()

		var deferred bool
		app.Defer(func() { deferred = true })

		var stop bool
		app.OnStop(func() { stop = true })

		go func() {
			time.Sleep(time.Second)
			syscall.Kill(syscall.Getpid(), sig)
		}()

		app.Run()

		if deferred == false {
			t.Error("Failed to invoke deferred functions.")
		}

		if stop == false {
			t.Error("Failed to invoke stop functions.")
		}
	}
}

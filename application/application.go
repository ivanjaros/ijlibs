package application

import (
	"os"
	"os/signal"
	"syscall"
)

type tasks []func()

func (t tasks) invoke() {
	if t != nil {
		for _, task := range t {
			if task != nil {
				task()
			}
		}
	}
}

type App interface {
	Run()
	Defer(tasks ...func()) App
	OnStop(tasks ...func()) App
	OnForcedStop(tasks ...func()) App
	Stop()
}

func New() App {
	return new(instance)
}

type instance struct {
	sigChan         chan os.Signal
	doneChan        chan struct{}
	stopChan        chan struct{}
	closing         bool
	forced          bool
	deferred        tasks
	stopTasks       tasks
	forcedStopTasks tasks
}

func (app *instance) Defer(tasks ...func()) App {
	app.deferred = append(app.deferred, tasks...)
	return app
}

func (app *instance) OnStop(tasks ...func()) App {
	app.stopTasks = append(app.stopTasks, tasks...)
	return app
}

func (app *instance) OnForcedStop(tasks ...func()) App {
	app.forcedStopTasks = append(app.forcedStopTasks, tasks...)
	return app
}

func (app *instance) Stop() {
	app.stopChan <- struct{}{}
}

func (app *instance) done() {
	select {
	case app.doneChan <- struct{}{}:
	default:
	}
}

func (app *instance) shutdown(force bool) {
	defer app.done()

	if app.forced {
		return
	}

	if app.closing || force {
		app.forced = true
		app.forcedStopTasks.invoke()
		return
	}

	app.closing = true
	app.stopTasks.invoke()
}

func (app *instance) Run() {
	app.sigChan = make(chan os.Signal, 1)
	app.doneChan = make(chan struct{}, 1)
	app.stopChan = make(chan struct{}, 1)
	defer app.deferred.invoke()

	// Info on signals: https://www.gnu.org/software/libc/manual/html_node/Termination-Signals.html
	signal.Notify(
		app.sigChan,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGKILL,
	)

	for {
		select {
		case s := <-app.sigChan:
			go app.shutdown(s == syscall.SIGKILL)

		case <-app.stopChan:
			go app.shutdown(app.forced)

		case <-app.doneChan:
			return
		}
	}
}

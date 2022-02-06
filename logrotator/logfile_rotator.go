package logrotator

import (
	"bytes"
	"github.com/ivanjaros/ijlibs/files"
	"io"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func New(filePath string) (r *LogFileRotator, e error) {
	r = &LogFileRotator{
		path:    filePath,
		closeCh: make(chan struct{}),
		buff:    new(bytes.Buffer),
		ingress: make(chan []byte),
		rotate:  make(chan string),
	}

	e = r.open()

	if e == nil {
		go func() {
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGHUP)

			for {
				select {
				case <-r.closeCh:
					return

				case msg := <-r.ingress:
					if r.handle == nil {
						r.buff.Write(msg)
					} else {
						if _, err := r.handle.Write(msg); err != nil {
							log.Println("failed to write log entry: " + err.Error())
						}
					}

				case <-sigChan:
					if err := r.close(); err != nil {
						log.Println("error closing log file: ", err)
					} else {
						if err := r.open(); err != nil {
							log.Println("error opening log file: ", err)
						}
					}

				case rotate := <-r.rotate:
					r.doRotate(rotate)
				}
			}
		}()
	}

	return
}

type LogFileRotator struct {
	path     string
	handle   *os.File
	closeCh  chan struct{}
	buff     *bytes.Buffer
	ingress  chan []byte
	rotate   chan string
	closer   sync.Once
	closeErr error
}

func (r *LogFileRotator) open() error {
	if err := r.close(); err != nil {
		return err
	}

	handle, err := os.OpenFile(r.path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	r.handle = handle

	if _, err := io.Copy(r.handle, r.buff); err != nil {
		return err
	}

	r.buff.Reset()

	return nil
}

func (r *LogFileRotator) close() error {
	if r.handle == nil {
		return nil
	}

	err := r.handle.Close()
	r.handle = nil

	return err
}

func (r *LogFileRotator) Close() error {
	r.closer.Do(func() {
		close(r.closeCh)
		close(r.ingress)
		r.closeErr = r.close()
	})
	return r.closeErr
}

// this always returns 0, nil
func (r *LogFileRotator) Write(p []byte) (n int, err error) {
	select {
	case <-r.closeCh:
	case r.ingress <- p:
	default:
	}
	return 0, nil
}

// calling this will cause the existing log file to be copied into the fileName and truncated to 0.
func (r *LogFileRotator) Rotate(fileName string) {
	r.rotate <- fileName
}

func (r *LogFileRotator) doRotate(fileName string) {
	// this will never happen but just in case
	if r.handle == nil {
		log.Println("cannot rotate log into " + fileName + " because the log file is inactive")
		return
	}

	if err := files.Copy(r.path, fileName); err != nil {
		log.Println("failed to rotate log into " + fileName + ": " + err.Error())
		return
	}

	if err := r.handle.Truncate(0); err != nil {
		log.Println("failed to truncate log file during rotation: " + err.Error())
		return
	}
	if _, err := r.handle.Seek(0, 0); err != nil {
		log.Println("failed to reset offset for active log file during rotation: " + err.Error())
		return
	}
}

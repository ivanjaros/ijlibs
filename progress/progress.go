package progress

import (
	"log"
	"os"
	"strings"
)

func New() *Writer {
	w := &Writer{writeCh: make(chan string), closeCh: make(chan struct{})}

	go func() {
		defer close(w.closeCh)
		var ln int
		for {
			select {
			case v, ok := <-w.writeCh:
				if ok == false {
					if _, err := os.Stdout.WriteString("\a"); err != nil {
						log.Print(err)
					}
					return
				}

				if len(v) < ln {
					v += strings.Repeat(" ", ln-len(v))
				}

				if _, err := os.Stdout.WriteString("\r" + v); err != nil {
					log.Print(err)
				}

				ln = len(v)
			}
		}
	}()

	return w
}

type Writer struct {
	writeCh chan string
	closeCh chan struct{}
}

func (p *Writer) Set(value string) {
	// non-blocking writes for fast producers
	select {
	case p.writeCh <- value:
	default:
	}
}

func (p *Writer) Stop() {
	close(p.writeCh)
	<-p.closeCh
}

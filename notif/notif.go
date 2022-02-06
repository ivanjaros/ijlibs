package notif

import (
	"context"
	"sync"
	"time"
	"unsafe"
)

type Client struct {
	recv   chan interface{}
	ch     *Channel
	o      sync.Once
	ctx    context.Context
	cancel context.CancelFunc
}

// will be nil if this client is write-only
func (c *Client) Listen() <-chan interface{} {
	return c.recv
}

func (c *Client) Close() {
	select {
	case <-c.ctx.Done():
	case c.ch.unsubscribe <- c:
	}
}

func (c *Client) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *Client) doClose() {
	c.o.Do(func() {
		c.cancel()
		if c.recv != nil {
			close(c.recv)
		}
	})
}

func (c *Client) send(msg interface{}) {
	// write-only clients will not handle any messages
	if c.recv == nil {
		return
	}
	t := time.NewTimer(c.ch.sc)
	select {
	case <-c.ctx.Done():
	case c.recv <- msg:
	case <-t.C:
		// time out/slow consumer, close the connection
		c.Close()
	}
}

func (c *Client) Broadcast(payload interface{}) bool {
	select {
	case <-c.ctx.Done():
		return false
	default:
		c.ch.Broadcast() <- &envelope{Message: payload, Sender: uintptr(unsafe.Pointer(c))}
		return true
	}
}

type envelope struct {
	Message interface{}
	Sender  uintptr
}

// leech is channel-blocking so goroutine should be called internally to make it non-blocking
// this is to ensure proper order of leeched messages.
func NewChannel(ctx context.Context, name string, slowConsumer time.Duration, emptyCh chan string, leech func(interface{})) *Channel {
	return &Channel{
		name:        name,
		ingres:      make(chan interface{}, 1000),
		subscribe:   make(chan *Client, 1000),
		unsubscribe: make(chan *Client, 1000),
		aud:         make(map[*Client]struct{}, 1000),
		ctx:         ctx,
		sc:          slowConsumer,
		empty:       emptyCh,
		leech:       leech,
		counter:     make(chan chan<- int, 10),
	}
}

type Channel struct {
	name        string
	ingres      chan interface{}
	subscribe   chan *Client
	unsubscribe chan *Client
	aud         map[*Client]struct{}
	ctx         context.Context
	sc          time.Duration
	empty       chan string
	leech       func(interface{})
	counter     chan chan<- int
}

func (ch *Channel) Id() string {
	return ch.name
}

// subscription is read-write by default. by providing "writeOnly=true", it can be switched into write-only mode
// in which case the client will not be disconnected for being slow reader.
func (ch *Channel) Subscribe(writeOnly ...bool) *Client {
	c := &Client{
		ch: ch,
	}
	if len(writeOnly) == 0 || writeOnly[0] == false {
		c.recv = make(chan interface{})
	}
	c.ctx, c.cancel = context.WithCancel(ch.ctx)
	ch.subscribe <- c
	return c
}

func (ch *Channel) Broadcast() chan<- interface{} {
	return ch.ingres
}

// returns once context is cancelled
func (ch *Channel) Start() {
	for {
		select {
		case <-ch.ctx.Done():
			for cl := range ch.aud {
				delete(ch.aud, cl)
				cl.doClose()
			}
			return
		case cl := <-ch.subscribe:
			ch.aud[cl] = struct{}{}

		case cl := <-ch.unsubscribe:
			delete(ch.aud, cl)
			cl.doClose()
			if len(ch.aud) == 0 {
				ch.signalEmpty()
			}

		case msg := <-ch.ingres:
			e, ok := msg.(*envelope)
			if ok {
				msg = e.Message
			}
			for cl := range ch.aud {
				if ok == false || uintptr(unsafe.Pointer(cl)) != e.Sender {
					go cl.send(e.Message)
				}
			}
			if ch.leech != nil {
				ch.leech(msg)
			}

		case count := <-ch.counter:
			count <- len(ch.aud)
		}
	}
}

// returns number of clients in this channel
func (ch *Channel) Count() int {
	req := make(chan int)
	ch.counter <- req
	n := <-req
	close(req)
	return n
}

// the same as Count() but it uses caller's channel
func (ch *Channel) CountChan(req chan int) {
	ch.counter <- req
}

func (ch *Channel) signalEmpty() {
	if ch.empty == nil {
		return
	}

	select {
	case ch.empty <- ch.name:
	default:
	}
}

type subscribeRequest struct {
	name string
	recv chan *Client
	wo   bool
}

type broadcastRequest struct {
	name string
	recv chan *Channel
}

type hasRequest struct {
	name string
	recv chan bool
}

type brokeredChannel struct {
	ch     *Channel
	cancel context.CancelFunc
}

type brokerLeech interface {
	Match(string) func(interface{})
}

func NewBroker(ctx context.Context, slowConsumer time.Duration, leech brokerLeech) *Broker {
	return &Broker{
		chans:     make(map[string]*brokeredChannel, 100),
		subscribe: make(chan *subscribeRequest, 10),
		broadcast: make(chan *broadcastRequest, 10),
		ctx:       ctx,
		sc:        slowConsumer,
		empty:     make(chan string, 10),
		leech:     leech,
		counter:   make(chan chan<- int, 10),
		has:       make(chan hasRequest, 10),
	}
}

type Broker struct {
	chans     map[string]*brokeredChannel
	subscribe chan *subscribeRequest
	broadcast chan *broadcastRequest
	ctx       context.Context
	sc        time.Duration
	empty     chan string
	leech     brokerLeech
	counter   chan chan<- int
	has       chan hasRequest
}

// returns once context is cancelled
func (b *Broker) Start() {
	for {
		select {
		case <-b.ctx.Done():
			return
		case req := <-b.subscribe:
			ch, ok := b.chans[req.name]
			if ok == false {
				ctx, cancel := context.WithCancel(b.ctx)
				var l func(interface{})
				if b.leech != nil {
					l = b.leech.Match(req.name)
				}
				ch = &brokeredChannel{
					ch:     NewChannel(ctx, req.name, b.sc, b.empty, l),
					cancel: cancel,
				}
				b.chans[req.name] = ch
				go ch.ch.Start()
			}
			req.recv <- ch.ch.Subscribe(req.wo)

		case req := <-b.broadcast:
			ch, ok := b.chans[req.name]
			if ok == false {
				ctx, cancel := context.WithCancel(b.ctx)
				var l func(interface{})
				if b.leech != nil {
					l = b.leech.Match(req.name)
				}
				ch = &brokeredChannel{
					ch:     NewChannel(ctx, req.name, b.sc, b.empty, l),
					cancel: cancel,
				}
				b.chans[req.name] = ch
				go ch.ch.Start()
			}
			req.recv <- ch.ch

		case name := <-b.empty:
			if ch, ok := b.chans[name]; ok {
				ch.cancel()
				delete(b.chans, name)
			}

		case count := <-b.counter:
			count <- len(b.chans)

		case has := <-b.has:
			_, ok := b.chans[has.name]
			has.recv <- ok
		}
	}
}

// subscription is read-write by default. by providing "writeOnly=true", it can be switched into write-only mode
// in which case the client will not be disconnected for being slow reader.
func (b *Broker) Subscribe(name string, writeOnly ...bool) *Client {
	req := &subscribeRequest{
		name: name,
		recv: make(chan *Client),
		wo:   len(writeOnly) > 0 && writeOnly[0] == true,
	}
	b.subscribe <- req
	c := <-req.recv
	close(req.recv)
	return c
}

func (b *Broker) Broadcast(name string) chan<- interface{} {
	req := &broadcastRequest{name: name, recv: make(chan *Channel)}
	b.broadcast <- req
	ch := <-req.recv
	close(req.recv)
	return ch.ingres
}

// returns number of channels in this broker
func (b *Broker) Count() int {
	req := make(chan int)
	b.counter <- req
	n := <-req
	close(req)
	return n
}

// the same as Count() but it uses caller's channel
func (b *Broker) CountChan(req chan int) {
	b.counter <- req
}

func (b *Broker) Has(name string) bool {
	req := hasRequest{name: name, recv: make(chan bool)}
	b.has <- req
	has := <-req.recv
	close(req.recv)
	return has
}

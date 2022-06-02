package events

type Subscriber chan Event

type Hub interface {
	Register(s Subscriber)
	Unregister(s Subscriber)
	Close()
	Publish(e Event)
}

type hub struct {
	ch          chan Event
	stop        chan struct{}
	register    chan Subscriber
	unregister  chan Subscriber
	subscribers map[Subscriber]bool
}

func NewHub() Hub {
	h := &hub{
		ch:          make(chan Event),
		stop:        make(chan struct{}),
		register:    make(chan Subscriber),
		unregister:  make(chan Subscriber),
		subscribers: make(map[Subscriber]bool),
	}

	go h.run()

	return h
}

func (h *hub) run() {
	for {
		select {
		case e := <-h.ch:
			for s := range h.subscribers {
				s <- e
			}
		case s := <-h.register:
			h.subscribers[s] = true
		case s := <-h.unregister:
			delete(h.subscribers, s)
		case <-h.stop:
			return
		}
	}
}

func (h *hub) Register(s Subscriber) {
	h.register <- s
}

func (h *hub) Unregister(s Subscriber) {
	h.unregister <- s
}

func (h *hub) Close() {
	close(h.stop)
}

func (h *hub) Publish(e Event) {
	h.ch <- e
}

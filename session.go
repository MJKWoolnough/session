package session

import (
	"sync"
	"time"
)

type object struct {
	data    interface{}
	timeout *time.Timer
}

func (o *object) run(s *Session, name string) {
	_, ok := <-o.timeout.C
	if ok {
		s.mu.Lock()
		defer s.mu.Unlock()
		ob, ok := s.cache[name]
		if ok && ob == o {
			delete(s.cache, name)
		}
	}
}

type Session struct {
	timeout time.Duration

	mu    sync.Mutex
	cache map[string]*object
}

func New(timeout time.Duration) *Session {
	return &Session{
		timeout: timeout,
		cache:   make(map[string]*object),
	}
}

func (s *Session) Set(name string, data interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if o, ok := s.cache[name]; ok {
		o.timeout.Stop()
	}
	o := &object{
		data:    data,
		timeout: time.NewTimer(s.timeout),
	}
	go o.run(s, name)
}

func (s *Session) Get(name string) interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	o, ok := s.cache[name]
	if !ok {
		return nil
	}
	o.timeout.Reset(s.timeout)
	return o.data
}

func (s *Session) Remove(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if o, ok := s.cache[name]; ok {
		o.timeout.Stop()
		delete(s.cache, name)
	}
}

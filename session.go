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

// Session stores data.
//
// The data will be removed from storage after a timeout has expired.
// The timeout is reset if the data is read from storage.
type Session struct {
	timeout time.Duration

	mu    sync.Mutex
	cache map[string]*object
}

// New creates a new session with the given timeout
func New(timeout time.Duration) *Session {
	return &Session{
		timeout: timeout,
		cache:   make(map[string]*object),
	}
}

// Set places the data into the session, setting the timer.
//
// If something with the same name already exists it will be removed from the
// storage
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

// Get retrieves a value from storage, resetting the timer when it does so
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

// Remove removes the named item from storage
func (s *Session) Remove(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if o, ok := s.cache[name]; ok {
		o.timeout.Stop()
		delete(s.cache, name)
	}
}

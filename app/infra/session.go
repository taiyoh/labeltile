package infra

import (
	"strings"
	"time"

	"github.com/taiyoh/labeltile/app"
)

type SessionData struct {
	app.SessionData
	data map[string]interface{}
}

func (d *SessionData) Get(name string) interface{} {
	o, _ := d.data[name]
	return o
}

func (d *SessionData) Set(name string, val interface{}) {
	d.data[name] = val
}

func (d *SessionData) Remove(name string) bool {
	if _, ok := d.data[name]; ok {
		delete(d.data, name)
		return true
	}
	return false
}

// MemorySession provides in-memory session storage for simple system
type MemorySession struct {
	app.SessionStorage
	storage map[string]app.SessionData
	prefix  string
	expire  uint32
}

func (s *MemorySession) sessionKey(id string) string {
	return strings.Join([]string{s.prefix, id}, ":")
}

func (s *MemorySession) Find(id string) app.SessionData {
	sid := s.sessionKey(id)
	if d, ok := s.storage[sid]; ok {
		return d
	}
	return nil
}

func (s *MemorySession) New(id string) app.SessionData {
	sid := s.sessionKey(id)
	if _, ok := s.storage[sid]; ok {
		return nil
	}
	d := &SessionData{data: map[string]interface{}{}}
	s.storage[sid] = d
	go func() {
		time.After(time.Duration(s.expire) * time.Second)
		s.Remove(id)
	}()
	return d
}

func (s *MemorySession) Save(id string, d app.SessionData) {
	s.storage[s.sessionKey(id)] = d
}

func (s *MemorySession) Remove(id string) bool {
	sid := s.sessionKey(id)
	if _, ok := s.storage[sid]; ok {
		delete(s.storage, sid)
		return true
	}
	return false
}

func NewMemorySession(prefix string, expire uint32) app.SessionStorage {
	return &MemorySession{
		storage: map[string]app.SessionData{},
		prefix:  prefix,
		expire:  expire,
	}
}

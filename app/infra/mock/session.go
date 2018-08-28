package mock

import (
	"strings"

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

type Session struct {
	app.SessionStorage
	storage map[string]app.SessionData
	prefix  string
	expire  uint32
}

func (s *Session) sessionKey(id string) string {
	return strings.Join([]string{s.prefix, id}, ":")
}

func (s *Session) Find(id string) app.SessionData {
	sid := s.sessionKey(id)
	if d, ok := s.storage[sid]; ok {
		return d
	}
	return nil
}

func (s *Session) New(id string) app.SessionData {
	sid := s.sessionKey(id)
	if _, ok := s.storage[sid]; ok {
		return nil
	}
	d := &SessionData{data: map[string]interface{}{}}
	s.storage[sid] = d
	return d
}

func (s *Session) Save(id string, d app.SessionData) {
	s.storage[s.sessionKey(id)] = d
}

func (s *Session) Remove(id string) bool {
	sid := s.sessionKey(id)
	if _, ok := s.storage[sid]; ok {
		delete(s.storage, sid)
		return true
	}
	return false
}

func LoadSession(prefix string, expire uint32) app.SessionStorage {
	return &Session{storage: map[string]app.SessionData{}, prefix: prefix, expire: expire}
}

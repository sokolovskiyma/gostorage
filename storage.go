package gostorage

import (
	"encoding/gob"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/sokolovskiyma/gostorage/item"
)

// type comerableCustom interface {
// 	Bytes() []byte
// }

// type stringCustom string

// func (s *stringCustom) Bytes() []byte {
// 	return []byte{}
// }

// type Number interface {
// 	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
// }

type storage[V any] struct {
	mu                 sync.RWMutex
	items              map[string]*item.Item[V]
	cleaner            *cleaner[V]
	fetch              func(string) (V, error)
	cleanupIntrval     time.Duration
	defalultExpiration int64
}

// SETUP

func (s *storage[V]) WithExpiration(defalultExpiration time.Duration) Storage[V] {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.defalultExpiration = int64(defalultExpiration)
	return s
}

func (s *storage[V]) WithCleaner(cleanupIntrval time.Duration) Storage[V] {
	s.mu.Lock()
	defer s.mu.Unlock()

	if cleanupIntrval > 0 {
		s.cleaner = &cleaner[V]{
			Interval: cleanupIntrval,
			stop:     make(chan bool),
		}
		go s.cleaner.Run(s)

		runtime.SetFinalizer(s, stopCleaner[V])
	}

	return s
}

// ACTIONS

func (s *storage[V]) SaveFile(filename string) (err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return
	}
	defer file.Close()

	err = gob.NewEncoder(file).Encode(s.items)
	if err != nil {
		return
	}

	return
}

func (s *storage[T]) LoadFile(filename string) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	err = gob.NewDecoder(file).Decode(&s.items)
	if err != nil {
		return
	}

	return
}

func (s *storage[V]) DeleteExpired() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixNano()
	for key, value := range s.items {
		if value.Expiration > 0 && now > value.Expiration {
			delete(s.items, key)
		}
	}
}

// FUNCTIONS

func (s *storage[V]) Set(key string, value V) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.set(key, value)
}

func (s *storage[V]) set(key string, value V) {
	if s.defalultExpiration == 0 {
		s.items[key] = &item.Item[V]{
			Value:      value,
			Expiration: 0,
		}
	} else {
		s.items[key] = &item.Item[V]{
			Value:      value,
			Expiration: time.Now().UnixNano() + s.defalultExpiration,
		}
	}
}

func (s *storage[V]) Get(key string) (V, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.get(key)
}

func (s *storage[V]) get(key string) (V, bool) {
	var defalultValue V

	if item, found := s.items[key]; found {
		if item.Expiration == 0 || item.Expiration > time.Now().UnixNano() {
			return item.Value, true
		}
	}

	if s.fetch != nil {
		value, err := s.fetch(key)
		if err == nil {
			s.set(key, value)
			return value, true
		}
	}

	return defalultValue, false
}

func (s *storage[V]) GetFetch(key string, f func(string) (V, error)) (V, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, ok := s.get(key)
	if ok {
		return value, ok
	}

	if value, err := f(key); err == nil {
		s.set(key, value)
		return value, true
	}

	return value, false
}

func (s *storage[V]) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.items, key)
}

func (s *storage[V]) Keys() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	var keys = make([]string, 0, len(s.items))
	for key := range s.items {
		keys = append(keys, key)
	}
	return keys
}

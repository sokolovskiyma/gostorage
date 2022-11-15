package gostorage

import (
	"encoding/gob"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/sokolovskiyma/gostorage/v2/item"
)

type storage[V any] struct {
	mu       *sync.RWMutex
	items    map[string]*item.Item[V]
	cleaner  *cleaner[V]
	settings Settings
}

func newStorage[V any](settings Settings) *storage[V] {
	storage := storage[V]{
		mu:       &sync.RWMutex{},
		items:    make(map[string]*item.Item[V]),
		settings: settings,
	}

	if storage.settings.CleanupInterval > 0 {
		storage.cleaner = &cleaner[V]{
			Interval: storage.settings.CleanupInterval,
			stop:     make(chan bool),
		}
		go storage.cleaner.Run(&storage)

		runtime.SetFinalizer(&storage, stopCleaner[V])
	}

	time.Now().IsZero()

	return &storage
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

	now := time.Now().Unix()
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
	if s.settings.Expiration == 0 {
		s.items[key] = &item.Item[V]{
			Value:      value,
			Expiration: 0,
		}
	} else {
		s.items[key] = &item.Item[V]{
			Value:      value,
			Expiration: time.Now().Unix() + s.settings.Expiration,
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
		if item.Expiration == 0 || item.Expiration > time.Now().Unix() {
			return item.Value, true
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

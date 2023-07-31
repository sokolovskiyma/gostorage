package gostorage

import (
	"runtime"
	"sync"
	"time"
)

type storage[V any] struct {
	mu       *sync.RWMutex
	items    map[string]item[V]
	cleaner  *cleaner[V]
	settings iternalSettings
}

type item[V any] struct {
	value      V
	expiration int64
}

func newStorage[V any](settings iternalSettings) *storage[V] {
	storage := storage[V]{
		mu:       &sync.RWMutex{},
		items:    make(map[string]item[V]),
		settings: settings,
	}

	if storage.settings.cleanup > 0 {
		storage.cleaner = &cleaner[V]{
			Interval: storage.settings.cleanup,
			stop:     make(chan bool),
		}
		go storage.cleaner.Run(&storage)

		runtime.SetFinalizer(&storage, stopCleaner[V])
	}

	return &storage
}

// ACTIONS

// func (s *storage[V]) SaveFile(filename string) (err error) {
// 	s.mu.RLock()
// 	defer s.mu.RUnlock()

// 	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
// 	if err != nil {
// 		return
// 	}
// 	defer file.Close()

// 	err = gob.NewEncoder(file).Encode(s.items)
// 	if err != nil {
// 		return
// 	}

// 	return
// }

// func (s *storage[T]) LoadFile(filename string) (err error) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return
// 	}
// 	defer file.Close()

// 	err = gob.NewDecoder(file).Decode(&s.items)
// 	if err != nil {
// 		return
// 	}

// 	return
// }

func (s *storage[V]) DeleteExpired() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixNano()
	for key, value := range s.items {
		if value.expiration > 0 && now > value.expiration {
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
	if s.settings.expiration == 0 {
		s.items[key] = item[V]{
			value:      value,
			expiration: 0,
		}
	} else {
		s.items[key] = item[V]{
			value:      value,
			expiration: time.Now().UnixNano() + s.settings.expiration,
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
		if item.expiration == 0 || item.expiration > time.Now().UnixNano() {
			return item.value, true
		}
	}

	return defalultValue, false
}

func (s *storage[V]) Fetch(key string, f func(string) (V, bool)) (V, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, ok := s.get(key)
	if ok {
		return value, ok
	}

	if value, ok := f(key); ok {
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

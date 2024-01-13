package gostorage

import (
	"runtime"
	"sync"
	"time"
)

type storage[K comparable, V any] struct {
	mu       *sync.RWMutex
	items    map[K]item[V]
	cleaner  *cleaner[K, V]
	settings iternalSettings
}

type item[V any] struct {
	value      V
	expiration int64
}

func newStorage[K comparable, V any](settings iternalSettings) *storage[K, V] {
	storage := storage[K, V]{
		mu:       &sync.RWMutex{},
		items:    make(map[K]item[V]),
		settings: settings,
	}

	if storage.settings.cleanup > 0 {
		storage.cleaner = &cleaner[K, V]{
			Interval: storage.settings.cleanup,
			stop:     make(chan bool),
		}
		go storage.cleaner.Run(&storage)

		runtime.SetFinalizer(&storage, stopCleaner[K, V])
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

func (s *storage[K, V]) DeleteExpired() {
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

func (s *storage[K, V]) Set(key K, value V) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.set(key, value)
}

func (s *storage[K, V]) set(key K, value V) {
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

func (s *storage[K, V]) Get(key K) (V, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.get(key)
}

func (s *storage[K, V]) get(key K) (V, bool) {
	if item, found := s.items[key]; found {
		if item.expiration == 0 || item.expiration > time.Now().UnixNano() {
			return item.value, true
		}
	}

	var defalultValue V
	return defalultValue, false
}

func (s *storage[K, V]) Fetch(key K, f func(K) (V, bool)) (V, bool) {
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

func (s *storage[K, V]) Delete(key K) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.items, key)
}

func (s *storage[K, V]) Keys() []K {
	s.mu.Lock()
	defer s.mu.Unlock()

	var keys = make([]K, 0, len(s.items))
	for key := range s.items {
		keys = append(keys, key)
	}
	return keys
}

package gostorage

import (
	"encoding/gob"
	"os"
	"runtime"
	"sync"
	"time"
)

type Storage[T any] struct {
	sync.RWMutex
	items              map[string]Item[T]
	cleaner            *cleaner[T]
	cleanupIntrval     time.Duration
	defalultExpiration time.Duration
}

type Item[T any] struct {
	Value      T
	Expiration int64
}

func NewStorage[T any]() *Storage[T] {
	return &Storage[T]{
		items:              make(map[string]Item[T]),
		cleanupIntrval:     0,
		defalultExpiration: 0,
	}
}

// Setup

func (s *Storage[T]) DefaultExpiration(defalultExpiration time.Duration) *Storage[T] {
	s.Lock()
	defer s.Unlock()
	s.defalultExpiration = defalultExpiration
	return s
}

func (s *Storage[T]) WithCleaner(cleanupIntrval time.Duration) *Storage[T] {
	s.Lock()
	defer s.Unlock()

	if cleanupIntrval > 0 {
		s.cleaner = &cleaner[T]{
			Interval: cleanupIntrval,
			stop:     make(chan bool),
		}
		go s.cleaner.Run(s)

		runtime.SetFinalizer(s, stopCleaner[T])
	}

	return s
}

// Actions

func (s *Storage[T]) SaveFile(filename string) (err error) {
	s.Lock()
	defer s.Unlock()

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

func (s *Storage[T]) LoadFile(filename string) (err error) {
	s.Lock()
	defer s.Unlock()

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

func (s *Storage[T]) DeleteExpired() {
	s.Lock()
	defer s.Unlock()
	var toDelete []string
	now := time.Now().UnixNano()

	for key, value := range s.items {
		if value.Expiration > 0 && now > value.Expiration {
			toDelete = append(toDelete, key)
		}
	}

	for _, key := range toDelete {
		delete(s.items, key)
	}
}

// func (s *Storage) FromMap(data map[string]interface{}) *Storage {
// 	s.Lock()
// 	defer s.Unlock()
// 	return s
// }
// func (s *Storage) FromFile(filename string) *Storage {
// 	s.Lock()
// 	defer s.Unlock()
// 	return s
// }

// func newCacheWithJanitor(de time.Duration, ci time.Duration, m map[string]Item) *Cache {
// 	c := newCache(de, m)
// 	// This trick ensures that the janitor goroutine (which--granted it
// 	// was enabled--is running DeleteExpired on c forever) does not keep
// 	// the returned C object from being garbage collected. When it is
// 	// garbage collected, the finalizer stops the janitor goroutine, after
// 	// which c can be collected.
// 	C := &Cache{c}
// 	if ci > 0 {
// 	  runJanitor(c, ci)
// 	  runtime.SetFinalizer(C, stopJanitor)
// 	}
// 	return C
//   }

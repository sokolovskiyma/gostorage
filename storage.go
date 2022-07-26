package gostorage

import (
	"encoding/gob"
	"os"
	"runtime"
	"sync"
	"time"
)

type Storage[T any] struct {
	mu                 sync.RWMutex
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
	s.mu.Lock()
	defer s.mu.Unlock()
	s.defalultExpiration = defalultExpiration
	return s
}

func (s *Storage[T]) WithCleaner(cleanupIntrval time.Duration) *Storage[T] {
	s.mu.Lock()
	defer s.mu.Unlock()

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

func (s *Storage[T]) LoadFile(filename string) (err error) {
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

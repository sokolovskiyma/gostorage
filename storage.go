package gostorage

import (
	"encoding/gob"
	"os"
	"runtime"
	"sync"
	"time"
)

type Storage struct {
	sync.RWMutex
	items              map[string]Item
	cleaner            *cleaner
	cleanupIntrval     time.Duration
	defalultExpiration time.Duration
}

func NewStorage() *Storage {
	return &Storage{
		items:              make(map[string]Item),
		cleanupIntrval:     0,
		defalultExpiration: 0,
	}
}

// Setup

func (s *Storage) DefaultExpiration(defalultExpiration time.Duration) *Storage {
	s.Lock()
	defer s.Unlock()
	s.defalultExpiration = defalultExpiration
	return s
}

func (s *Storage) WithCleaner(cleanupIntrval time.Duration) *Storage {
	s.Lock()
	defer s.Unlock()

	if cleanupIntrval > 0 {
		s.cleaner = &cleaner{
			Interval: cleanupIntrval,
			stop:     make(chan bool),
		}
		go s.cleaner.Run(s)

		runtime.SetFinalizer(s, stopCleaner)
	}

	return s
}

// Actions

func (s *Storage) SaveFile(filename string) (err error) {
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

func (s *Storage) LoadFile(filename string) (err error) {
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

func (s *Storage) DeleteExpired() {
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

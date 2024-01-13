package gostorage

import "time"

type cleaner[K comparable, V any] struct {
	Interval time.Duration
	stop     chan bool
}

func (c *cleaner[K, V]) Run(s *storage[K, V]) {
	ticker := time.NewTicker(c.Interval)
	for {
		select {
		case <-ticker.C:
			s.DeleteExpired()
		case <-c.stop:
			ticker.Stop()
			return
		}
	}
}

func stopCleaner[K comparable, V any](s *storage[K, V]) {
	s.cleaner.stop <- true
}

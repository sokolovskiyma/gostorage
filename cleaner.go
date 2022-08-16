package gostorage

import "time"

type cleaner[V any] struct {
	Interval time.Duration
	stop     chan bool
}

func (c *cleaner[V]) Run(s *storage[V]) {
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

func stopCleaner[V any](s *storage[V]) {
	s.cleaner.stop <- true
}

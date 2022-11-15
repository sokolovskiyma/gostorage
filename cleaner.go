package gostorage

import "time"

type cleaner[V any] struct {
	Interval int64
	stop     chan bool
}

func (c *cleaner[V]) Run(s *storage[V]) {
	ticker := time.NewTicker(time.Duration(c.Interval) * time.Second)
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

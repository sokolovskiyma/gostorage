package gostorage

import "time"

type cleaner[T any] struct {
	Interval time.Duration
	stop     chan bool
}

func (c *cleaner[T]) Run(s *Storage[T]) {
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

func stopCleaner[T any](s *Storage[T]) {
	s.cleaner.stop <- true
}

//   func runJanitor(c *cache, ci time.Duration) {
// 	j := &janitor{
// 	  Interval: ci,
// 	  stop:     make(chan bool),
// 	}
// 	c.janitor = j
// 	go j.Run(c)
//   }

package gostorage

import "time"

type cleaner struct {
	Interval time.Duration
	stop     chan bool
}

func (c *cleaner) Run(s *Storage) {
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

func stopCleaner(s *Storage) {
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

package gostorage

import "time"

func (s *Storage[T]) Set(key string, value T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.defalultExpiration == 0 {
		s.items[key] = Item[T]{
			Value:      value,
			Expiration: 0,
		}
	} else {
		s.items[key] = Item[T]{
			Value:      value,
			Expiration: time.Now().Add(s.defalultExpiration).UnixNano(),
		}
	}

}

func (s *Storage[T]) SetForever(key string, value T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[key] = Item[T]{
		Value:      value,
		Expiration: 0,
	}
}

func (s *Storage[T]) SetTemporarily(key string, value T, expiration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[key] = Item[T]{
		Value:      value,
		Expiration: time.Now().Add(expiration).UnixNano(),
	}
}

func (s *Storage[T]) Get(key string) (T, bool) {
	var defalultValue T

	s.mu.RLock()
	defer s.mu.RUnlock()

	item, found := s.items[key]
	if !found {
		return defalultValue, false
	}
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			return defalultValue, false
		}
	}

	return item.Value, true
}

func (s *Storage[T]) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.items, key)
}

func (s *Storage[T]) DeleteExpired() {
	s.mu.Lock()
	defer s.mu.Unlock()
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

// func (s *Storage) Update(key string, value interface{}) error {
// 	s.Lock()
// 	defer s.Unlock()
// 	// TODO
// 	return nil
// }

// func (s *Storage) UpdateOrSet(key string, value interface{}) {
// 	s.Lock()
// 	defer s.Unlock()
// 	// TODO
// }

// func (s *Storage) UpdateWithExpiration(key string, value interface{}, expiration time.Duration) error {
// 	s.Lock()
// 	defer s.Unlock()
// 	// TODO
// 	return nil
// }

// func (s *Storage) UpdateOrSetWithExpiration(key string, value interface{}) {
// 	s.Lock()
// 	defer s.Unlock()
// 	// TODO
// }

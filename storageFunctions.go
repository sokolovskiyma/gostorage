package storage

import "time"

func (s *Storage) Set(key string, value interface{}) {
	s.Lock()
	defer s.Unlock()

	if s.defalultExpiration == 0 {
		s.items[key] = Item{
			Value:      value,
			Expiration: 0,
		}
	} else {
		s.items[key] = Item{
			Value:      value,
			Expiration: time.Now().Add(s.defalultExpiration).UnixNano(),
		}
	}

}

func (s *Storage) SetWithExpiration(key string, value interface{}, expiration time.Duration) {
	s.Lock()
	defer s.Unlock()

	s.items[key] = Item{
		Value:      value,
		Expiration: time.Now().Add(expiration).UnixNano(),
	}
}

func (s *Storage) Get(key string) (interface{}, bool) {
	s.Lock()
	defer s.Unlock()

	item, found := s.items[key]
	if !found {
		return nil, false
	}
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			return nil, false
		}
	}

	return item.Value, true
}

func (s *Storage) Delete(key string) {
	s.Lock()
	defer s.Unlock()

	delete(s.items, key)
}

// func (s *Storage) SetDefault(key string, value interface{}) {
// 	s.Lock()
// 	defer s.Unlock()

// 	var e int64
// 	if s.defalultExpiration != 0 {
// 		e = time.Now().Add(s.defalultExpiration).UnixNano()
// 	}

// 	s.items[key] = Item{
// 		Value:      value,
// 		Expiration: e,
// 	}
// }

// func (s *Storage) GetWithExpiration(key string) (interface{}, time.Duration, bool) {
// 	s.Lock()
// 	defer s.Unlock()

// 	// "Inlining" of get and Expired
// 	item, found := s.items[key]
// 	if !found {
// 		return nil, 0, false
// 	}
// 	if item.Expiration > 0 {
// 		if time.Now().UnixNano() > item.Expiration {
// 			return nil, 0, false
// 		}
// 	}

// 	return item.Value, time.Duration(item.Expiration), true
// }

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

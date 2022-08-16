package gostorage

import (
	"time"

	"github.com/sokolovskiyma/gostorage/item"
)

type Storage[V any] interface {
	WithCleaner(time.Duration) Storage[V]
	WithExpiration(time.Duration) Storage[V]

	SaveFile(string) error
	LoadFile(string) error
	DeleteExpired()

	Get(string) (V, bool)
	GetFetch(string, func(string) (V, error)) (V, bool)
	Set(string, V)
	Delete(string)
	Keys() []string
}

func NewStorage[V any]() Storage[V] {
	return &storage[V]{
		items:              make(map[string]*item.Item[V]),
		fetch:              nil,
		cleanupIntrval:     0,
		defalultExpiration: 0,
	}
}

func NewStorageShards[V any](numShards uint64) Storage[V] {
	var ss storageShards[V]
	for i := 0; i < int(numShards); i++ {
		// stor := NewStorage[T]()
		ss.shards = append(ss.shards, &storage[V]{
			items:              make(map[string]*item.Item[V]),
			cleanupIntrval:     0,
			defalultExpiration: 0,
		})
	}
	return &ss
}

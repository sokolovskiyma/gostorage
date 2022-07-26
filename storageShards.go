package gostorage

import (
	"encoding/gob"
	"os"
	"sync"
	"time"
)

type StorageShards[T any] struct {
	mu     sync.RWMutex
	shards []*Storage[T]
}

func NewStorageShards[T any](numShards int) *StorageShards[T] {
	if numShards < 1 {
		panic("numShards must be >= 1")
	}
	var ss StorageShards[T]
	for i := 0; i < numShards; i++ {
		// stor := NewStorage[T]()
		ss.shards = append(ss.shards, NewStorage[T]())
	}
	return &ss
}

// Setup

func (ss *StorageShards[T]) DefaultExpiration(defalultExpiration time.Duration) *StorageShards[T] {
	for index := range ss.shards {
		ss.shards[index].DefaultExpiration(defalultExpiration)
	}
	return ss
}

func (ss *StorageShards[T]) WithCleaner(cleanupIntrval time.Duration) *StorageShards[T] {
	for index := range ss.shards {
		ss.shards[index].WithCleaner(cleanupIntrval)
	}
	return ss
}

// Actions

func (ss *StorageShards[T]) SaveFile(filename string) (err error) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return
	}
	defer file.Close()

	// ERROR: Нельзя просто так сохранить shards
	// тк там есть пол которые не экспортируются
	// нужно создать новый make([]map[K]V, len(ss.shards))
	// сохранить туда осколки в цикле
	temp := make([]map[string]Item[T], 0, len(ss.shards))
	for index := range ss.shards {
		temp = append(temp, ss.shards[index].items)
	}
	err = gob.NewEncoder(file).Encode(temp)
	if err != nil {
		return
	}

	return
}

func (ss *StorageShards[T]) LoadFile(filename string) (err error) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	temp := make([]map[string]Item[T], 0, len(ss.shards))
	err = gob.NewDecoder(file).Decode(&temp)
	if err != nil {
		return
	}

	for index := range temp {
		ss.shards[index].items = temp[index]
	}

	return
}

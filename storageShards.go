package gostorage

import (
	"encoding/gob"
	"hash/fnv"
	"os"
	"sync"
	"time"

	"github.com/sokolovskiyma/gostorage/v2/item"
)

type storageShards[V any] struct {
	mu     sync.RWMutex
	shards []*storage[V]
}

// Setup

func (ss *storageShards[V]) WithExpiration(defalultExpiration time.Duration) Storage[V] {
	for index := range ss.shards {
		ss.shards[index].WithExpiration(defalultExpiration)
	}
	return ss
}

func (ss *storageShards[V]) WithCleaner(cleanupIntrval time.Duration) Storage[V] {
	for index := range ss.shards {
		ss.shards[index].WithCleaner(cleanupIntrval)
	}
	return ss
}

// Actions

func (ss *storageShards[V]) SaveFile(filename string) (err error) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return
	}
	defer file.Close()

	temp := make([]map[string]*item.Item[V], 0, len(ss.shards))
	for index := range ss.shards {
		temp = append(temp, ss.shards[index].items)
	}
	err = gob.NewEncoder(file).Encode(temp)
	if err != nil {
		return
	}

	return
}

func (ss *storageShards[V]) LoadFile(filename string) (err error) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	temp := make([]map[string]*item.Item[V], 0, len(ss.shards))
	err = gob.NewDecoder(file).Decode(&temp)
	if err != nil {
		return
	}

	for index := range temp {
		ss.shards[index].items = temp[index]
	}

	return
}

func (ss *storageShards[V]) DeleteExpired() {
	for index := range ss.shards {
		ss.shards[index].DeleteExpired()
	}
}

// FUNCTIONS

func (ss *storageShards[V]) shardByKey(key string) *storage[V] {
	hash := fnv.New32()
	hash.Write([]byte(key))
	return ss.shards[int(hash.Sum32())%len(ss.shards)]
}

func (ss *storageShards[V]) Set(key string, value V) {
	ss.shardByKey(key).Set(key, value)
}

func (ss *storageShards[V]) Get(key string) (V, bool) {
	return ss.shardByKey(key).Get(key)
}

func (ss *storageShards[V]) GetFetch(key string, f func(string) (V, error)) (V, bool) {
	return ss.shardByKey(key).GetFetch(key, f)
}

func (ss *storageShards[V]) Delete(key string) {
	ss.shardByKey(key).Delete(key)
}

func (ss *storageShards[V]) Keys() []string {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	var keys = make([]string, 0, 512)
	for shardIndex := range ss.shards {
		keys = append(keys, ss.shards[shardIndex].Keys()...)
	}

	return keys
}

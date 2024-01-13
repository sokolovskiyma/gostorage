package gostorage

import (
	"hash/maphash"
)

type storageShards[K comparable, V any] struct {
	shards []*storage[K, V]
	mask   uint64
	seed   maphash.Seed
	ksize  int
	kstr   bool
}

// Setup

func newStorageShards[K comparable, V any](settings iternalSettings) *storageShards[K, V] {
	storage := storageShards[K, V]{
		mask: uint64(settings.shards - 1),
		seed: maphash.MakeSeed(),
	}

	storage.detectHasher()

	for i := 0; i < int(settings.shards); i++ {
		storage.shards = append(storage.shards, newStorage[K, V](settings))
	}

	return &storage
}

// Actions

// func (ss *storageShards[V]) SaveFile(filename string) (err error) {
// 	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
// 	if err != nil {
// 		return
// 	}
// 	defer file.Close()

// 	temp := make([]map[string]*Item[V], 0, len(ss.shards))
// 	for index := range ss.shards {
// 		temp = append(temp, ss.shards[index].items)
// 	}
// 	err = gob.NewEncoder(file).Encode(temp)
// 	if err != nil {
// 		return
// 	}

// 	return
// }

// func (ss *storageShards[V]) LoadFile(filename string) (err error) {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return
// 	}
// 	defer file.Close()

// 	temp := make([]map[string]*Item[V], 0, len(ss.shards))
// 	err = gob.NewDecoder(file).Decode(&temp)
// 	if err != nil {
// 		return
// 	}

// 	for index := range temp {
// 		ss.shards[index].items = temp[index]
// 	}

// 	return
// }

func (ss *storageShards[K, V]) DeleteExpired() {
	for index := range ss.shards {
		ss.shards[index].DeleteExpired()
	}
}

// FUNCTIONS

func (ss *storageShards[K, V]) shardByKey(key K) *storage[K, V] {
	// hash := fnv.New64()
	// _, _ = hash.Write([]byte(key))
	// return ss.shards[int(hash.Sum64()&ss.mask)]

	// TODO: maphash
	return ss.shards[ss.hash(key)&ss.mask]
}

func (ss *storageShards[K, V]) Set(key K, value V) {
	ss.shardByKey(key).Set(key, value)
}

func (ss *storageShards[K, V]) Get(key K) (V, bool) {
	return ss.shardByKey(key).Get(key)
}

func (ss *storageShards[K, V]) Fetch(key K, f func(K) (V, bool)) (V, bool) {
	return ss.shardByKey(key).Fetch(key, f)
}

func (ss *storageShards[K, V]) Delete(key K) {
	ss.shardByKey(key).Delete(key)
}

func (ss *storageShards[K, V]) Keys() []K {
	var keys = make([]K, 0, 512)
	for shardIndex := range ss.shards {
		keys = append(keys, ss.shards[shardIndex].Keys()...)
	}

	return keys
}

package gostorage

import (
	"hash/maphash"
)

type storageShards[V any] struct {
	shards []*storage[V]
	mask   uint64
	seed   maphash.Seed
}

// Setup

func newStorageShards[V any](settings iternalSettings) *storageShards[V] {
	storage := storageShards[V]{
		mask: uint64(settings.shards - 1),
		seed: maphash.MakeSeed(),
	}

	for i := 0; i < int(settings.shards); i++ {
		storage.shards = append(storage.shards, newStorage[V](settings))
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

func (ss *storageShards[V]) DeleteExpired() {
	for index := range ss.shards {
		ss.shards[index].DeleteExpired()
	}
}

// FUNCTIONS

func (ss *storageShards[V]) shardByKey(key string) *storage[V] {
	// hash := fnv.New64()
	// _, _ = hash.Write([]byte(key))
	// return ss.shards[int(hash.Sum64()&ss.mask)]

	// TODO: maphash
	return ss.shards[maphash.String(ss.seed, key)&ss.mask]
}

func (ss *storageShards[V]) Set(key string, value V) {
	ss.shardByKey(key).Set(key, value)
}

func (ss *storageShards[V]) Get(key string) (V, bool) {
	return ss.shardByKey(key).Get(key)
}

func (ss *storageShards[V]) Fetch(key string, f func(string) (V, bool)) (V, bool) {
	return ss.shardByKey(key).Fetch(key, f)
}

func (ss *storageShards[V]) Delete(key string) {
	ss.shardByKey(key).Delete(key)
}

func (ss *storageShards[V]) Keys() []string {
	var keys = make([]string, 0, 512)
	for shardIndex := range ss.shards {
		keys = append(keys, ss.shards[shardIndex].Keys()...)
	}

	return keys
}

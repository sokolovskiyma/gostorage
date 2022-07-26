package gostorage

import (
	"hash/fnv"
	"time"
)

func (ss *StorageShards[T]) shardByKey(key string) *Storage[T] {
	hash := fnv.New32()
	hash.Write([]byte(key))
	return ss.shards[int(hash.Sum32())%len(ss.shards)]
}

func (ss *StorageShards[T]) Set(key string, value T) {
	ss.shardByKey(key).Set(key, value)
}

func (ss *StorageShards[T]) SetForever(key string, value T) {
	ss.shardByKey(key).SetForever(key, value)
}

func (ss *StorageShards[T]) SetTemporarily(key string, value T, expiration time.Duration) {
	ss.shardByKey(key).SetTemporarily(key, value, expiration)
}

func (ss *StorageShards[T]) Get(key string) (T, bool) {
	return ss.shardByKey(key).Get(key)
}

func (ss *StorageShards[T]) Delete(key string) {
	ss.shardByKey(key).Delete(key)
}

func (ss *StorageShards[T]) DeleteExpired() {
	for index := range ss.shards {
		ss.shards[index].DeleteExpired()
	}
}

package gostorage

import "time"

type Storage[K comparable, V any] interface {
	// SaveFile(string) error
	// LoadFile(string) error
	DeleteExpired()

	Get(K) (V, bool)
	Fetch(K, func(K) (V, bool)) (V, bool)
	Set(K, V)
	Delete(K)
	Keys() []K
}

type (
	Settings struct {
		Expiration time.Duration
		Cleanup    time.Duration
		Shards     uint32
	}
	iternalSettings struct {
		expiration int64
		cleanup    time.Duration
		shards     int
	}
)

func EmptySettings() Settings {
	return Settings{
		Expiration: 0,
		Cleanup:    0,
		Shards:     1,
	}
}

func DefaultSettings(expiration time.Duration) Settings {
	return Settings{
		Expiration: expiration,
		Cleanup:    0,
		Shards:     1,
	}
}

func NewStorage[K comparable, V any](settings Settings) Storage[K, V] {
	if settings.Shards == 0 {
		settings.Shards = 1
	}

	set := iternalSettings{
		expiration: int64(settings.Expiration),
		cleanup:    settings.Cleanup,
		shards:     int(settings.Shards),
	}

	if settings.Shards > 1 {
		return newStorageShards[K, V](set)
	}

	return newStorage[K, V](set)
}

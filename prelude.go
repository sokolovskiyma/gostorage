package gostorage

import "time"

type Storage[V any] interface {
	// SaveFile(string) error
	// LoadFile(string) error
	DeleteExpired()

	Get(string) (V, bool)
	Fetch(string, func(string) (V, bool)) (V, bool)
	Set(string, V)
	Delete(string)
	Keys() []string
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

func DefalultSettings(expiration time.Duration) Settings {
	return Settings{
		Expiration: expiration,
		Cleanup:    0,
		Shards:     1,
	}
}

func NewStorage[V any](settings Settings) Storage[V] {
	if settings.Shards == 0 {
		settings.Shards = 1
	}

	// TODO: int keys
	// func (h *hasher[K]) detectHasher() {
	// 	var k K
	// 	switch ((interface{})(k)).(type) {
	// 	case string:
	// 		h.kstr = true
	// 	default:
	// 		h.ksize = int(unsafe.Sizeof(k))
	// 	}
	// }

	set := iternalSettings{
		expiration: int64(settings.Expiration),
		cleanup:    settings.Cleanup,
		shards:     int(settings.Shards),
	}

	if settings.Shards > 1 {
		return newStorageShards[V](set)
	}

	return newStorage[V](set)
}

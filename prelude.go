package gostorage

type Storage[V any] interface {
	// WithCleaner(time.Duration) Storage[V]
	// WithExpiration(time.Duration) Storage[V]

	SaveFile(string) error
	LoadFile(string) error
	DeleteExpired()

	Get(string) (V, bool)
	GetFetch(string, func(string) (V, error)) (V, bool)
	Set(string, V)
	Delete(string)
	Keys() []string
}

type Settings struct {
	Expiration      int64
	CleanupInterval int64
	ShardsQuantity  uint64
}

func EmptySettings() Settings {
	return Settings{
		Expiration:      0,
		CleanupInterval: 0,
		ShardsQuantity:  1,
	}
}

func DefalultSettings(expiration int64) Settings {
	return Settings{
		Expiration:      int64(expiration),
		CleanupInterval: 0,
		ShardsQuantity:  1,
	}
}

func NewStorage[V any](settings Settings) Storage[V] {
	if settings.ShardsQuantity > 1 {
		return newStorageShards[V](settings)
	}

	return newStorage[V](settings)

}

package gostorage

type Item[V any] struct {
	Value      V
	Expiration int64
}

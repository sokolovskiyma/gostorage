package gostorage

import (
	"hash/maphash"
	"unsafe"
)

func (ss *storageShards[K, V]) detectHasher() {
	// Detect the key type. This is needed by the hasher.
	var k K
	switch ((interface{})(k)).(type) {
	case string:
		ss.kstr = true
	default:
		ss.ksize = int(unsafe.Sizeof(k))
	}
}

func (ss *storageShards[K, V]) hash(key K) uint64 {
	// The unsafe package is used here to cast the key into a string container
	// so that the hasher can work. The hasher normally only accept a string or
	// []byte, but this effectively allows it to accept value type.
	// The m.kstr bool, which is set from the New function, indicates that the
	// key is known to already be a true string. Otherwise, a fake string is
	// derived by setting the string data to value of the key, and the string
	// length to the size of the value.
	var strKey string
	if ss.kstr {
		strKey = *(*string)(unsafe.Pointer(&key))
	} else {
		strKey = *(*string)(unsafe.Pointer(&struct {
			data unsafe.Pointer
			len  int
		}{unsafe.Pointer(&key), ss.ksize}))
	}
	// Now for the actual hashing.

	return maphash.String(ss.seed, strKey)
}

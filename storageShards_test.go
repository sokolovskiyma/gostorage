package gostorage

import (
	"strconv"
	"testing"
	"time"
)

/* TESTS */

func TestNewStorageShards(t *testing.T) {
	stor := NewStorage[any, any](Settings{
		Expiration: 0,
		Cleanup:    0,
		Shards:     5,
	})

	if stor == nil {
		t.Log("stor == nil")
		t.Fail()
	}
}

func TestSetShards(t *testing.T) {

	// preparation
	stor := NewStorage[string, string](Settings{
		Expiration: 0,
		Cleanup:    0,
		Shards:     5,
	})

	// test
	stor.Set(testKey, testValue)
}

func TestGetShards(t *testing.T) {

	// preparation
	stor := NewStorage[string, string](Settings{
		Expiration: 0,
		Cleanup:    0,
		Shards:     5,
	})

	// test
	stor.Set(testKey, testValue)
	if value, ok := stor.Get(testKey); !ok {
		t.Log("there is no value 'test'")
		t.Fail()
	} else if value != testValue {
		t.Logf("value %+v != %+v\n", value, testValue)
		t.Fail()
	}

	// test
	if value, ok := stor.Get("nonexist"); ok || value != "" {
		t.Log("found nonexistent value")
		t.Fail()
	}
}

func TestWithFetchShards(t *testing.T) {

	// preparation
	stor := NewStorage[string, string](Settings{
		Expiration: 0,
		Cleanup:    0,
		Shards:     5,
	})

	// test
	value, ok := stor.Fetch(testKey, func(s string) (string, bool) {
		return testValue, true
	})

	if !ok {
		t.Log("there is no value 'test'")
		t.Fail()
	}

	if value != testValue {
		t.Logf("value %+v != %+v\n", value, testValue)
		t.Fail()
	}
}

func TestDeleteShards(t *testing.T) {

	// preparation
	stor := NewStorage[string, string](Settings{
		Expiration: 0,
		Cleanup:    0,
		Shards:     5,
	})

	// test
	stor.Set(testKey, testValue)
	stor.Delete(testKey)

	if value, ok := stor.Get(testKey); ok || value != "" {
		t.Log("found deleted value")
		t.Fail()
	}
}

func TestKeysShards(t *testing.T) {

	// preparation
	stor := NewStorage[string, string](Settings{
		Expiration: 0,
		Cleanup:    0,
		Shards:     5,
	})
	stor.Set("one", testValue)
	stor.Set("two", testValue)
	stor.Set("three", testValue)

	// test
	if len(stor.Keys()) != 3 {
		t.Log("missing some keys")
		t.Fail()
	}

	for _, key := range stor.Keys() {
		if _, ok := stor.Get(key); !ok {
			t.Logf("there is no key %v", key)
			t.Fail()
		}
	}
}

func TestGetWithExpirationShards(t *testing.T) {
	// preparation
	stor := NewStorage[string, string](Settings{
		Expiration: 0,
		Cleanup:    0,
		Shards:     5,
	})

	// test
	stor.Set(testKey, testValue)

	if value, ok := stor.Get(testKey); !ok {
		t.Log("there is no value 'test'")
		t.Fail()
	} else if value != testValue {
		t.Logf("value %+v != %+v\n", value, testValue)
		t.Fail()
	}

	stor = NewStorage[string, string](Settings{
		Expiration: 5,
		Cleanup:    0,
		Shards:     5,
	})
	stor.Set(testKey, testValue)
	time.Sleep(testSleep)

	if value, ok := stor.Get(testKey); ok || value != "" {
		t.Log("found expired value")
		t.Fail()
	}
}

func TestWithExpirationShards(t *testing.T) {
	// preparation
	stor := NewStorage[string, string](Settings{
		Expiration: testExpiration,
		Cleanup:    0,
		Shards:     5,
	})

	// test
	stor.Set(testKey, testValue)

	if value, ok := stor.Get(testKey); !ok {
		t.Log("there is no value 'test'")
		t.Fail()
	} else if value != testValue {
		t.Logf("value %+v != %+v\n", value, testValue)
		t.Fail()
	}

	time.Sleep(testSleep)

	if value, ok := stor.Get(testKey); ok || value != "" {
		t.Log("found expired value")
		t.Fail()
	}
}

func TestWithCleanerShards(t *testing.T) {

	// preparation
	stor := NewStorage[string, string](Settings{
		Expiration: 0,
		Cleanup:    1,
		Shards:     5,
	})

	// test
	stor.Set(testKey, testValue)
	time.Sleep(testSleep)

	if value, ok := stor.Get(testKey); !ok || value == "" {
		t.Log("there is no value 'test'")
		t.Fail()
	} else if value != testValue {
		t.Logf("value %+v != %+v\n", value, testValue)
		t.Fail()
	}

	// test
	stor = NewStorage[string, string](Settings{
		Expiration: 5,
		Cleanup:    1,
		Shards:     5,
	})
	stor.Set(testKey, testValue)
	time.Sleep(testSleep)

	if value, ok := stor.Get(testKey); ok || value != "" {
		t.Log("found expired value")
		t.Fail()
	}

	if _, ok := stor.(*storageShards[string, string]).shardByKey(testValue).items[testKey]; ok {
		t.Log("found deleted value")
		t.Fail()
	}
}

// func TestSaveLoadFileShards(t *testing.T) {

// 	// preparation
// 	stor := NewStorage[string, string](Settings{
// 		Expiration: testExpiration,
// 		Cleanup:    0,
// 		Shards:     5,
// 	})

// 	// test
// 	stor.Set(testKey, testValue)

// 	err := stor.SaveFile("testfile")
// 	if err != nil {
// 		t.Log("cant save file", err.Error())
// 		t.Fail()
// 	}

// 	// test
// 	stor2 := NewStorage[string, string](Settings{
// 		Expiration: 0,
// 		Cleanup:    0,
// 		Shards:     5,
// 	})
// 	err = stor2.LoadFile("testfile")
// 	if err != nil {
// 		t.Log("cant load file", err.Error())
// 		t.Fail()
// 	}

// 	if value, ok := stor2.Get(testKey); !ok {
// 		t.Log("there is no value 'test'")
// 		t.Fail()
// 	} else if value != testValue {
// 		t.Logf("value %+v != %+v\n", value, testValue)
// 		t.Fail()
// 	}

// }

/* BENCHMARKS */

func BenchmarkSetGetShards(b *testing.B) {
	stor := NewStorage[string, string](Settings{
		Expiration: 0,
		Cleanup:    0,
		Shards:     32,
	})

	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(b.N)
		stor.Set(key, key)
		value, has := stor.Get(key)
		if !has {
			b.Fatal("!has")
		}
		if value != key {
			b.Fatal("value != key")
		}
	}
}

func BenchmarkSetGetWithExpirationShards(b *testing.B) {
	stor := NewStorage[string, string](Settings{
		Expiration: 5,
		Cleanup:    0,
		Shards:     32,
	})

	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(b.N)
		stor.Set(key, key)
		stor.Get(key)
	}
}

func BenchmarkSetGetSetWithFetchShards(b *testing.B) {
	stor := NewStorage[string, string](Settings{
		Expiration: 0,
		Cleanup:    0,
		Shards:     32,
	})

	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(b.N)
		stor.Fetch(key, func(k string) (string, bool) {
			return k, true
		})
	}
}

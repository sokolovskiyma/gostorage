package gostorage

import (
	"testing"
	"time"
)

// TODO
// const (
// 	testKeyShards        = "testKeyShards"
// 	testValueShards      = "testValueShards"
// 	testExpirationShards = 5 * time.Second
// )

/* TESTS */

func TestNewStorageShards(t *testing.T) {
	stor := NewStorageShards[any](5)

	if stor == nil {
		t.Log("stor == nil")
		t.Fail()
	}
}

func TestSetShards(t *testing.T) {

	// preparation
	stor := NewStorageShards[string](5)

	// test
	stor.Set(testKey, testValue)
}

func TestGetShards(t *testing.T) {

	// preparation
	stor := NewStorageShards[string](5)

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
	stor := NewStorageShards[string](5)

	// test
	value, ok := stor.GetFetch(testKey, func(s string) (string, error) {
		return testValue, nil
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
	stor := NewStorageShards[string](5)

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
	stor := NewStorageShards[string](5)
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
	stor := NewStorageShards[string](5)

	// test
	stor.Set(testKey, testValue)

	if value, ok := stor.Get(testKey); !ok {
		t.Log("there is no value 'test'")
		t.Fail()
	} else if value != testValue {
		t.Logf("value %+v != %+v\n", value, testValue)
		t.Fail()
	}

	stor.WithExpiration(testExpiration)
	stor.Set(testKey, testValue)
	time.Sleep(testExpiration)

	if value, ok := stor.Get(testKey); ok || value != "" {
		t.Log("found expired value")
		t.Fail()
	}
}

func TestWithExpirationShards(t *testing.T) {
	// preparation
	stor := NewStorageShards[string](5).WithExpiration(testExpiration)

	// test
	stor.Set(testKey, testValue)

	if value, ok := stor.Get(testKey); !ok {
		t.Log("there is no value 'test'")
		t.Fail()
	} else if value != testValue {
		t.Logf("value %+v != %+v\n", value, testValue)
		t.Fail()
	}

	time.Sleep(testExpiration)

	if value, ok := stor.Get(testKey); ok || value != "" {
		t.Log("found expired value")
		t.Fail()
	}
}

func TestWithCleanerShards(t *testing.T) {

	// preparation
	stor := NewStorageShards[string](5).WithCleaner(time.Second)

	// test
	stor.Set(testKey, testValue)
	time.Sleep(testExpiration)

	if value, ok := stor.Get(testKey); !ok || value == "" {
		t.Log("there is no value 'test'")
		t.Fail()
	} else if value != testValue {
		t.Logf("value %+v != %+v\n", value, testValue)
		t.Fail()
	}

	// test
	stor.WithExpiration(testExpiration)
	stor.Set(testKey, testValue)
	time.Sleep(testExpiration)

	if value, ok := stor.Get(testKey); ok || value != "" {
		t.Log("found expired value")
		t.Fail()
	}

	if _, ok := stor.(*storageShards[string]).shardByKey(testValue).items[testKey]; ok {
		t.Log("found deleted value")
		t.Fail()
	}
}

func TestSaveLoadFileShards(t *testing.T) {

	// preparation
	stor := NewStorageShards[string](5).WithExpiration(testExpiration)

	// test
	stor.Set(testKey, testValue)

	err := stor.SaveFile("testfile")
	if err != nil {
		t.Log("cant save file", err.Error())
		t.Fail()
	}

	// test
	stor2 := NewStorageShards[string](5)
	err = stor2.LoadFile("testfile")
	if err != nil {
		t.Log("cant load file", err.Error())
		t.Fail()
	}

	if value, ok := stor2.Get(testKey); !ok {
		t.Log("there is no value 'test'")
		t.Fail()
	} else if value != testValue {
		t.Logf("value %+v != %+v\n", value, testValue)
		t.Fail()
	}

}

/* BENCHMARKS */

func BenchmarkSetGetShards(b *testing.B) {
	// preparation
	stor := NewStorageShards[string](5)

	// test
	for i := 0; i < b.N; i++ {
		stor.Set(testKey, testValue)
		stor.Get(testKey)
	}
}

func BenchmarkSetGetWithExpirationShards(b *testing.B) {
	// preparation
	stor := NewStorageShards[string](5).WithExpiration(testExpiration)

	// test
	for i := 0; i < b.N; i++ {
		stor.Set(testKey, testValue)
		stor.Get(testKey)
	}
}

func BenchmarkSetGetSetWithFetchShards(b *testing.B) {
	// preparation
	stor := NewStorageShards[string](5)

	// test
	for i := 0; i < b.N; i++ {
		stor.GetFetch(testKey, func(string) (string, error) {
			return testValue, nil
		})
	}
}

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

func TestGetWithExpirationShards(t *testing.T) {
	// preparation
	stor := NewStorageShards[string](5)

	// test
	stor.SetTemporarily(testKey, testValue, testExpiration)

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

func TestDefaultExpirationShards(t *testing.T) {
	// preparation
	stor := NewStorageShards[string](5).DefaultExpiration(testExpiration)

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

func TestCleanerShards(t *testing.T) {

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
	stor.SetTemporarily(testKey, testValue, time.Second)
	time.Sleep(testExpiration)

	if value, ok := stor.Get(testKey); ok || value != "" {
		t.Log("found expired value")
		t.Fail()
	}

	if _, ok := stor.shardByKey(testValue).items[testKey]; ok {
		t.Log("found deleted value")
		t.Fail()
	}
}

func TestSaveLoadFileShards(t *testing.T) {

	// preparation
	stor := NewStorageShards[string](5)

	// test
	stor.SetTemporarily(testKey, testValue, testExpiration)

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

func BenchmarkSetForeverShards(b *testing.B) {
	// preparation
	stor := NewStorageShards[string](5)

	// test
	for i := 0; i < b.N; i++ {
		stor.SetForever(testKey, testValue)
	}
}

func BenchmarkSetShards(b *testing.B) {
	// preparation
	stor := NewStorageShards[string](5)

	// test
	for i := 0; i < b.N; i++ {
		stor.Set(testKey, testValue)
	}
}

func BenchmarkSetDefaultShards(b *testing.B) {
	// preparation
	stor := NewStorageShards[string](5).DefaultExpiration(testExpiration)

	// test
	for i := 0; i < b.N; i++ {
		stor.Set(testKey, testValue)
	}
}

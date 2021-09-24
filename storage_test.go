package storage

import (
	"testing"
	"time"
)

const (
	testKey        = "testKey"
	testValue      = "testValue"
	testExpiration = 5 * time.Second
)

/* TESTS */

func TestNewStorage(t *testing.T) {
	stor := NewStorage()

	if stor == nil {
		t.Log("stor == nil")
		t.Fail()
	}
}

func TestSet(t *testing.T) {

	// preparation
	stor := NewStorage()

	// test
	stor.Set(testKey, testValue)
}

func TestGet(t *testing.T) {

	// preparation
	stor := NewStorage()

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
	if value, ok := stor.Get("nonexist"); ok || value != nil {
		t.Log("found nonexistent value")
		t.Fail()
	}
}

func TestDelete(t *testing.T) {

	// preparation
	stor := NewStorage()

	// test
	stor.Set(testKey, testValue)
	stor.Delete(testKey)

	if value, ok := stor.Get(testKey); ok || value != nil {
		t.Log("found deleted value")
		t.Fail()
	}
}

func TestGetWithExpiration(t *testing.T) {
	// preparation
	stor := NewStorage()

	// test
	stor.SetWithExpiration(testKey, testValue, testExpiration)

	if value, ok := stor.Get(testKey); !ok {
		t.Log("there is no value 'test'")
		t.Fail()
	} else if value != testValue {
		t.Logf("value %+v != %+v\n", value, testValue)
		t.Fail()
	}

	time.Sleep(testExpiration)

	if value, ok := stor.Get(testKey); ok || value != nil {
		t.Log("found expired value")
		t.Fail()
	}
}

func TestDefaultExpiration(t *testing.T) {
	// preparation
	stor := NewStorage().DefaultExpiration(testExpiration)

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

	if value, ok := stor.Get(testKey); ok || value != nil {
		t.Log("found expired value")
		t.Fail()
	}
}

func TestCleaner(t *testing.T) {

	// preparation
	stor := NewStorage().WithCleaner(testExpiration)

	// test
	stor.Set(testKey, testValue)
	time.Sleep(testExpiration)

	if value, ok := stor.Get(testKey); !ok || value == nil {
		t.Log("there is no value 'test'")
		t.Fail()
	} else if value != testValue {
		t.Logf("value %+v != %+v\n", value, testValue)
		t.Fail()
	}

	// test
	stor.SetWithExpiration(testKey, testValue, testExpiration)
	time.Sleep(testExpiration)

	if value, ok := stor.Get(testKey); ok || value != nil {
		t.Log("found expired value")
		t.Fail()
	}

	if _, ok := stor.items[testValue]; ok {
		t.Log("found deleted value")
		t.Fail()
	}
}

func TestSaveLoadFile(t *testing.T) {

	// preparation
	stor := NewStorage()

	// test
	stor.SetWithExpiration(testKey, testValue, testExpiration)

	err := stor.SaveFile("testfile")
	if err != nil {
		t.Log("cant save file", err.Error())
		t.Fail()
	}

	// test
	stor2 := NewStorage()
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

func BenchmarkSet(b *testing.B) {
	// preparation
	stor := NewStorage()

	// test
	for i := 0; i < b.N; i++ {
		stor.Set(testKey, testValue)
	}
}

func BenchmarkSetDefault(b *testing.B) {
	// preparation
	stor := NewStorage().DefaultExpiration(testExpiration)

	// test
	for i := 0; i < b.N; i++ {
		stor.Set(testKey, testValue)
	}
}

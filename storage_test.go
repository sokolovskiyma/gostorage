package gostorage

import (
	"errors"
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
	stor := NewStorage[any]()

	if stor == nil {
		t.Log("stor == nil")
		t.Fail()
	}
}

func TestSet(t *testing.T) {

	// preparation
	stor := NewStorage[string]()

	// test
	stor.Set(testKey, testValue)
}

func TestGet(t *testing.T) {

	// preparation
	stor := NewStorage[string]()

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

func TestWithFetch(t *testing.T) {
	// preparation
	stor := NewStorage[string]()

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

	_, ok = stor.GetFetch("foo", func(s string) (string, error) {
		return "", errors.New("no value")
	})

	if ok {
		t.Log("found nonexistent value")
		t.Fail()
	}
}

func TestDelete(t *testing.T) {

	// preparation
	stor := NewStorage[string]()

	// test
	stor.Set(testKey, testValue)
	stor.Delete(testKey)

	if value, ok := stor.Get(testKey); ok || value != "" {
		t.Log("found deleted value")
		t.Fail()
	}
}

func TestKeys(t *testing.T) {

	// preparation
	stor := NewStorage[string]()
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

func TestGetWithExpiration(t *testing.T) {
	// preparation
	stor := NewStorage[string]().WithExpiration(testExpiration)

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

func TestWithExpiration(t *testing.T) {
	// preparation
	stor := NewStorage[string]().WithExpiration(testExpiration)

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

func TestWithCleaner(t *testing.T) {

	// preparation
	stor := NewStorage[string]().WithCleaner(time.Second)

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
	stor.WithExpiration(time.Second)
	stor.Set(testKey, testValue)
	time.Sleep(testExpiration)

	if value, ok := stor.Get(testKey); ok || value != "" {
		t.Log("found expired value")
		t.Fail()
	}

	if _, ok := stor.(*storage[string]).items[testKey]; ok {
		t.Log("found deleted value")
		t.Fail()
	}
}

func TestSaveLoadFile(t *testing.T) {

	// preparation
	stor := NewStorage[string]().WithExpiration(testExpiration)

	// test
	stor.Set(testKey, testValue)

	err := stor.SaveFile("testfile")
	if err != nil {
		t.Log("cant save file", err.Error())
		t.Fail()
	}

	// test
	stor2 := NewStorage[string]()
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

func BenchmarkSetGet(b *testing.B) {
	// preparation
	stor := NewStorage[string]()

	// test
	for i := 0; i < b.N; i++ {
		stor.Set(testKey, testValue)
		stor.Get(testKey)
	}
}

func BenchmarkSetGetWithExpiration(b *testing.B) {
	// preparation
	stor := NewStorage[string]().WithExpiration(testExpiration)

	// test
	for i := 0; i < b.N; i++ {
		stor.Set(testKey, testValue)
		stor.Get(testKey)
	}
}

func BenchmarkSetGetSetWithFetch(b *testing.B) {
	// preparation
	stor := NewStorage[string]()

	// test
	for i := 0; i < b.N; i++ {
		stor.GetFetch(testKey, func(string) (string, error) {
			return testValue, nil
		})
	}
}

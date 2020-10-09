package main

import (
	"fmt"
	"strconv"
	"sync"
)

type (
	Key   = string
	Value = string
)

type Cache interface {
	GetOrSet(key Key, valueFn func() Value) Value
	Get(key Key) (Value, bool)
}

// ----------------------------------------------

type InMemoryCache struct {

	// Тип добавлен в версии Golang 1.9 https://golang.org/src/sync/map.go
	// Испоьзуется для случаев:
	// (1) когда запись ключа происходит только один раз,
	// но читается много раз, как в кэшах, которые только растут.
	// (2) когда несколько goroutines читают, записывают и перезаписывают значения
	// для непересекающихся наборов ключей.
	// В этих 2-х случаях использование sync.Mutex может значительно
	// уменьшить конкуренцию блокировок по сравнению с использование Mutex or RWMutex
	syncMutex sync.Mutex

	dataMutex sync.RWMutex
	data map[Key]Value
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		data: make(map[Key]Value),
	}
}

func (cache *InMemoryCache) Get(key Key) (Value, bool) {
	cache.dataMutex.RLock()
	defer cache.dataMutex.RUnlock()

	value, found := cache.data[key]
	return value, found
}

// GetOrSet возвращает значение ключа в случае его существования.
// Иначе, вычисляет значение ключа при помощи valueFn,
// сохраняет его в кэш и возвращает это значение.
func (cache *InMemoryCache) GetOrSet(key Key, valueFn func() Value) Value {

	cache.syncMutex.Lock()
	defer cache.syncMutex.Unlock()

	if value, found := cache.data[key]; found {
		return value
	}

	val := valueFn()

	cache.data[key] = val

	return val
}

// Метод для тестирования производительности
// Сравнение sync.Mutex с sync.RWMutex
func (cache *InMemoryCache) GetOrSetRWMutex(key Key, valueFn func() Value) Value {

	value, found := cache.Get(key)

	if found {
		return value
	}

	cache.dataMutex.Lock()
	defer cache.dataMutex.Unlock()

	val := valueFn()

	cache.data[key] = val

	return val
}

func main() {

	cache := NewInMemoryCache()

	i := 0

	fn := func() Value {
		i++
		return strconv.Itoa(i)
	}

	var wg sync.WaitGroup

	for count := 0; count <= 10000; count++ {

		wg.Add(1)

		go func() {
			cache.GetOrSet("some-key-1", fn)
			wg.Done()
		}()

	}

	wg.Wait()

	fmt.Printf("Значение i после завершения 10000 вызовов cache.GetOrSet(\"some-key-1\", fn): %v\n", i)
}
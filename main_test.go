package main

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"strconv"
	"sync"
	"testing"
)

// 1. Тестирование того, что конкурентные обращения
// к существующиму ключу не блокирует друг друга
//
// 2. Тестирование того, что значение каждого ключа
// будет вычислено ровно 1 раз
func TestInMemoryCache_GetOrSet_fn_only_once(t *testing.T) {

	cache := NewInMemoryCache()

	key := uuid.New().String()

	i := 0

	fn := func() Value {
		i++
		return strconv.Itoa(i)
	}

	var wg sync.WaitGroup

	for count := 0; count <= 100000; count++ {

		wg.Add(1)

		go func() {
			cache.GetOrSet(key, fn)
			wg.Done()
		}()

	}

	wg.Wait()

	valFromCache := fn()

	assert.EqualValues(
		t,
		"2",
		valFromCache,
		"Функция func() Value была вызвана более, чем 1 раз или не вызывалась ни разу",
	)

}

// Тестирование различных значений key и values
func TestInMemoryCache_GetOrSet_values(t *testing.T) {

	cache := NewInMemoryCache()

	testCases := []struct{
		name string
		key Key
		fn func() Value
		resultValue string
		isEqual bool
	}{
		{
			name: "Equal key and value Chinese",
			key: "key-中文单词",
			fn: func() Value {
				return "中文单词"
			},
			resultValue: "中文单词",
			isEqual: true,
		},
		{
			name: "Not Equal string",
			key: "key-2",
			fn: func() Value {
				return "שפה גרמנית"
			},
			resultValue: "Election Live Updates: After Virus-Stricken Trump Rejects Virtual Debate",
			isEqual: false,
		},
		{
			name: "Not Equal char",
			key: "key-3",
			fn: func() Value {
				return "x"
			},
			resultValue: "X",
			isEqual: false,
		},
		{
			name: "Equal. Empty val",
			key: "key-4-empty-val",
			fn: func() Value {
				return ""
			},
			resultValue: "",
			isEqual: true,
		},
		{
			name: "Equal. Empty key",
			key: "",
			fn: func() Value {
				return "some value"
			},
			resultValue: "some value",
			isEqual: true,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {

			if tc.isEqual {
				assert.EqualValues(t, tc.resultValue, cache.GetOrSet(tc.key, tc.fn))
			} else {
				assert.NotEqualValues(t, tc.resultValue, cache.GetOrSet(tc.key, tc.fn))
			}

		})

	}

}

func BenchmarkInMemoryCache_GetOrSet(b *testing.B) {

	fn := func() Value {
		return "result string"
	}

	cache := NewInMemoryCache()

	for i := 0; i < b.N; i++ {
		for count := 0; count <= 100000; count++ {
			go func() {
				cache.GetOrSet("some-key", fn)
			}()
		}
	}

}

func BenchmarkInMemoryCache_GetOrSetRWMutex(b *testing.B) {

	fn := func() Value {
		return "result string"
	}

	cache := NewInMemoryCache()

	for i := 0; i < b.N; i++ {

		for count := 0; count <= 100000; count++ {

			go func() {
				cache.GetOrSetRWMutex("some-key", fn)
			}()

		}
	}

}

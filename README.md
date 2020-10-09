# InMemoryCache

InMemoryCache — потоко-безопасная реализация Key-Value кэша, хранящая данные в оперативной памяти.

## Использование

```go
...

type (
	Key   = string
	Value = string
)

func main() {

    cache := NewInMemoryCache()

    fn := func() Value {
    	return "some value"
    }

    // Устанавливает значение для ключа "key"
    val = cache.GetOrSet("key", fn)
    
    // Получает значение ключа "key", 
    // которое установлено раннее
    val = cache.GetOrSet("key", fn)
}
```

## Тестирование

Запуск тестов:

```shell script
make test
```

Результат:

```
go test -v -race ./
=== RUN   TestInMemoryCache_GetOrSet_fn_only_once
--- PASS: TestInMemoryCache_GetOrSet_fn_only_once (4.27s)
=== RUN   TestInMemoryCache_GetOrSet_values
=== RUN   TestInMemoryCache_GetOrSet_values/Equal_key_and_value_Chinese
=== RUN   TestInMemoryCache_GetOrSet_values/Not_Equal_string
=== RUN   TestInMemoryCache_GetOrSet_values/Not_Equal_char
=== RUN   TestInMemoryCache_GetOrSet_values/Equal._Empty_val
=== RUN   TestInMemoryCache_GetOrSet_values/Equal._Empty_key
--- PASS: TestInMemoryCache_GetOrSet_values (0.00s)
    --- PASS: TestInMemoryCache_GetOrSet_values/Equal_key_and_value_Chinese (0.00s)
    --- PASS: TestInMemoryCache_GetOrSet_values/Not_Equal_string (0.00s)
    --- PASS: TestInMemoryCache_GetOrSet_values/Not_Equal_char (0.00s)
    --- PASS: TestInMemoryCache_GetOrSet_values/Equal._Empty_val (0.00s)
    --- PASS: TestInMemoryCache_GetOrSet_values/Equal._Empty_key (0.00s)
PASS
ok      _/Users/drandin/code/in-memory-cache    4.566s

```

Запуск бенчмарков:

```shell script
make bench
```

Сравнивается работа методов **GetOrSet** и **GetOrSetRWMutex**. 

Перывый испоьзует **sync.Mutex**, второй — **sync.RWMutex**.

Результат:

```
go test -bench . -benchmem . ./
goos: darwin
goarch: amd64
BenchmarkInMemoryCache_GetOrSet-6                     60          18879793 ns/op           18841 B/op         48 allocs/op
BenchmarkInMemoryCache_GetOrSetRWMutex-6              48          23221211 ns/op              54 B/op          0 allocs/op
PASS
ok      _/Users/drandin/code/in-memory-cache    2.577s

```
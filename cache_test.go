package objcache

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

type ForTest struct {
	value string
}

func BenchmarkSimple(b *testing.B) {
	myconfig := Config{
		MaxEntryLimit: 200000,
		Expiration:    5 * time.Minute,
	}
	cache, err := New(myconfig)
	if err != nil {
		panic(err)
	}

	a := ForTest{
		value: "a",
	}

	//fmt.Println(time.Now().UnixNano())

	cache.Set("1", &a, 10*time.Minute)
	cache.Set("2", &a, 10*time.Minute)
	cache.Set("3", &a, 10*time.Minute)
	cache.Set("4", &a, 10*time.Minute)
	cache.Set("1", &a, 10*time.Minute)
	cache.Set("5", &a, 10*time.Minute)
	cache.Set("6", &a, 10*time.Minute)
	cache.Set("7", &a, 10*time.Minute)
	cache.Set("8", &a, 10*time.Minute)
	cache.Set("9", &a, 10*time.Minute)
	cache.Set("10", &a, 10*time.Minute)
	cache.Set("11", &a, 10*time.Minute)

	k := make(chan int)

	for i := 1; i <= 11; i = i + 1 {
		go func(xx int) {
			for j := 1; j < 1000000; j = j + 1 {
				x := strconv.Itoa(xx)
				_, ok := cache.Get(x)
				aa := ForTest{
					value: strconv.Itoa(j*11 + i),
				}
				cache.Set(strconv.Itoa(j*11+i), &aa, 10*time.Minute)
				if !ok {
					continue
				}

			}
			k <- 1
		}(i)
	}

	for j := 0; j < 11; j = j + 1 {
		c := <-k
		fmt.Println(c)
	}
	//fmt.Println(time.Now().UnixNano())
}

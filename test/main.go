package main

import (
	"fmt"
	"time"

	"strconv"

	".."
)

type ForTest struct {
	value string
}

func main() {
	myconfig := objcache.Config{
		MaxEntryLimit: 10,
		Expiration:    5 * time.Minute,
	}
	cache, err := objcache.New(myconfig)
	if err != nil {
		panic(err)
	}

	a := ForTest{
		value: "a",
	}

	cache.Set("1", &a, 10*time.Minute)
	cache.Set("2", &a, 10*time.Minute)
	cache.Set("3", &a, 10*time.Minute)
	cache.Set("4", &a, 10*time.Minute)
	cache.Set("5", &a, 10*time.Minute)
	cache.Set("6", &a, 10*time.Minute)
	cache.Set("7", &a, 10*time.Minute)
	cache.Set("8", &a, 10*time.Minute)
	cache.Set("9", &a, 10*time.Minute)
	cache.Set("10", &a, 10*time.Minute)
	cache.Set("11", &a, 10*time.Minute)

	for i := 1; i <= 11; i = i + 1 {
		x := strconv.Itoa(i)
		a, ok := cache.Get(x)

		if ok {
			b := a.(*ForTest)
			fmt.Println(x, b.value)
		} else {
			fmt.Println("FUCk")
		}
	}
}

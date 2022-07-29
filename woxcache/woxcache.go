package woxcache

import (
	"fmt"
	"time"
)

/**
This package exists to keep value pair data in cache for a specified TTL value.
It's a good practice to cache database results or files contents.
*/

const ttl = 3

type sCacheValueNVT struct {
	name  string
	value string
	ttl   int
}
type sCacheValueVT struct {
	value      string
	expiration time.Time
	//	ttl   int
}

var cacheEntries = make(map[string]sCacheValueVT)

func CacheInit() {
	//cacheEntries = make(map[string]sCacheValueVT)
}

func CachePurge() {
	purge := ""
	t := time.Now()
	//	fmt.Println("temps : ",t)
	for k, v := range cacheEntries {
		//
		//		fmt.Println("temps exp : ",v.expiration)
		temp := v.expiration.Sub(t)

		//.Second()
		//		fmt.Println("Temp: ",temp.Seconds() )
		//		fmt.Printf("Comparing ith %s",k)

		if temp.Seconds() < 0 {
			purge = k
			delete(cacheEntries, purge)
		}

	}
	fmt.Println("--> should purge : ", purge)
}

/**
Dump the entire cache to the screen for debugging
*/
func CacheDump() {
	if len(cacheEntries) == 0 {
		fmt.Println("-> Cache dump is empty")
	} else {
		fmt.Println("-> Cache dump have values: ")
		for k, v := range cacheEntries {
			fmt.Printf("   - Item [%s] value [%s] expiration [%v]\n", k, v.value, v.expiration)
			//			fmt.Println(v.value)
			//fmt.Println(v.ttl)
		}
	}
}

/**
store a value in the cache
*/
func CacheSet(name string, val string) string {
	t := time.Now()
	fmt.Println("-> Cache store [", name, "] ")
	temp := sCacheValueVT{value: val, expiration: t.Add(time.Second * ttl)}
	cacheEntries[name] = temp
	return val
}

/**
fetch a value from the cache
*/
func CacheGet(name string) string {
	for k, v := range cacheEntries {
		//		fmt.Printf("Comparing %s with %s",name,k)
		if k == name {
			return v.value
		}
	}
	return ""
}

package runtime_test

import (
	"fmt"
	"math"
	"time"
	"runtime"
	"testing"
)

// go clean -testcache ; time go test -v runtime/map_bench_m_64_64_test.go

func run(mapCacheFriendly bool) {
	runtime.MapCacheFriendly = mapCacheFriendly
	var keys int64 = 50000000
	var debug bool = false
	if debug {
		runtime.MapMakeDebug = true
		runtime.MapIterDebug = true
		keys = 10
	}
	m := make(map[int64]int64, keys*2) // *2 to not trigger overflow
	var i int64

	{
		var ir int64 = 0
		t1 := time.Now()
		for i = 1; i <= keys; i++ {
			m[i] = i
			ir += i
		}
		t2 := time.Now()
		elapsed := t2.Sub(t1).Seconds()
		fmt.Printf("- put %d map keys in %6.3f seconds or %10.0f keys per second; ir=%d mapCacheFriendly=%v\n", keys, elapsed, math.Floor(float64(keys) / elapsed), ir, mapCacheFriendly)
	}

	{
		var ir int64 = 0
		t1 := time.Now()
		for i = 1; i <= keys; i++ {
			ir += m[i]
		}
		t2 := time.Now()
		elapsed := t2.Sub(t1).Seconds()
		fmt.Printf("- got %d map vals in %6.3f seconds or %10.0f keys per second; ir=%d mapCacheFriendly=%v\n", keys, elapsed, math.Floor(float64(keys) / elapsed), ir, mapCacheFriendly)
	}

	if debug {
		for k, v := range m { 
			fmt.Printf(" - key[%+v] value[%+v]\n", k, v)
		}
	}
}

func TestMapPerformance(t *testing.T) {
	run(false); run(true)
	run(false); run(true)
	run(false); run(true)
	run(false); run(true)
}

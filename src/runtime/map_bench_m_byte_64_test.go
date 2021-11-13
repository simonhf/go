package runtime_test

import (
	"fmt"
	"math"
	"os"
	"time"
	"runtime"
    "runtime/debug"
	"testing"
)

// export MAP_BENCH_SCRIPT="map_bench_m_byte_64_test.go" ; export MAP_BENCH_DEBUG=true  ; cd ../src ; go clean -testcache ; time go test -v runtime/$MAP_BENCH_SCRIPT
// export MAP_BENCH_SCRIPT="map_bench_m_byte_64_test.go" ; export MAP_BENCH_DEBUG=false ; cd ../src ; go clean -testcache ; time go test -v runtime/$MAP_BENCH_SCRIPT | tee ../../$MAP_BENCH_SCRIPT.txt
// export MAP_BENCH_SCRIPT="map_bench_m_byte_64_test.go" ; cat ../../$MAP_BENCH_SCRIPT.txt | perl -lane 'if(m~(got|put|range).*mapCacheFriendly=(true|false)~){ push @{$h->{$1}{$2}}, $_; } sub END{ foreach $pg(qw(put got range)){ foreach $ft(qw(false true)){ @a=@{$h->{$pg}{$ft}}; foreach(@a){ if(m~(\d+) per sec~){ $r->{$pg}{$ft}+=$1; } printf qq[%s\n], $_; } printf qq[- %d total keys per second\n], $r->{$pg}{$ft}; } printf qq[- %.1f%% diff; true better than false if %% > 0\n], ($r->{$pg}{true} - $r->{$pg}{false}) / $r->{$pg}{false} * 100; } }'

var mapBenchDebug bool = os.Getenv("MAP_BENCH_DEBUG") != "false";

func run(mapCacheFriendly bool, verbose bool) {
	defer runtime.GC()
	runtime.MapCacheFriendly = mapCacheFriendly
	var keys int64 = 10000000
	if mapBenchDebug {
		runtime.MapMakeDebug = true
		runtime.MapIterDebug = true
		keys = 10
	}
	m := make(map[[32]byte]int64, keys*2) // *2 to not trigger overflow
	var i int64
	var irt int64 = 0
	key := [32]byte{0,0,0,0,0,0,0,0, 0,0,0,0,0,0,0,0, 0,0,0,0,0,0,0,0, 0,0,0,0,0,0,0,0}

	{
		var ir int64 = 0
		t1 := time.Now()
		for i = 1; i <= keys; i++ {
			key[0] = byte((i >>  0) & 255)
			key[1] = byte((i >>  8) & 255)
			key[2] = byte((i >> 16) & 255)
			key[3] = byte((i >> 24) & 255)
			m[key] = i
			ir += i
		}
		t2 := time.Now()
		elapsed := t2.Sub(t1).Seconds()
		if verbose {
			fmt.Printf("-   put %d map keys in %6.3f seconds or %10.0f per second; ir=%d mapCacheFriendly=%v\n", keys, elapsed, math.Floor(float64(keys) / elapsed), ir, mapCacheFriendly)
		}
		irt = ir
	}

	{
		var ir int64 = 0
		t1 := time.Now()
		for i = 1; i <= keys; i++ {
			key[0] = byte((i >>  0) & 255)
			key[1] = byte((i >>  8) & 255)
			key[2] = byte((i >> 16) & 255)
			key[3] = byte((i >> 24) & 255)
			ir += m[key]
		}
		t2 := time.Now()
		elapsed := t2.Sub(t1).Seconds()
		if verbose {
			fmt.Printf("-   got %d map vals in %6.3f seconds or %10.0f per second; ir=%d mapCacheFriendly=%v\n", keys, elapsed, math.Floor(float64(keys) / elapsed), ir, mapCacheFriendly)
		}
		if ir != irt { panic("ERROR: ir != 1st ir") }
	}

	{
		var ir int64 = 0
		t1 := time.Now()
		if mapBenchDebug {
			for k, v := range m { 
				ir += v
				fmt.Printf(" - key[%+v] value[%+v]\n", k, v)
			}
		} else {
			for _, v := range m { 
				ir += v
			}
		}
		t2 := time.Now()
		elapsed := t2.Sub(t1).Seconds()
		if verbose {
			fmt.Printf("- range %d map vals in %6.3f seconds or %10.0f per second; ir=%d mapCacheFriendly=%v\n", keys, elapsed, math.Floor(float64(keys) / elapsed), ir, mapCacheFriendly)
		}
		if ir != irt { panic("ERROR: ir != 1st ir") }
	}
}

func TestMapPerformance(t *testing.T) {
    defer debug.SetGCPercent(debug.SetGCPercent(-1)) // disable GC
	if !mapBenchDebug {
		run(false, false); // first run is always slower due to heap being expanded
	}
	run(false, true); run(true, true)
	if !mapBenchDebug {
		run(false, true); run(true, true)
		run(false, true); run(true, true)
		run(false, true); run(true, true)
	}
}

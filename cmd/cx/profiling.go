package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	cxcore "github.com/skycoin/cx/cx"
)

var DebugProfile bool
var DebugProfileRate int

var profiles = make(map[string]int64)

func StartProfile(name string) {
	if DebugProfile {
		profiles[name] = time.Now().UnixNano()
	}
}

func StopProfile(name string) {
	if DebugProfile {
		t := time.Now().UnixNano()
		deltaTime := t - profiles[name]
		fmt.Printf("%s : %dms\n", name, deltaTime/(int64(time.Millisecond)))
	}
}

func StartCPUProfile(name string) *os.File {
	if DebugProfile {
		f, err := cxcore.CXCreateFile(fmt.Sprintf("%s_%s_cpu.pprof", os.Args[0], name))
		if err != nil {
			fmt.Println("Failed to create CPU profile: ", err)
		}
		if DebugProfileRate != 100 { // test against default value to avoid warning
			runtime.SetCPUProfileRate(DebugProfileRate)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			fmt.Println("Failed to start CPU profile: ", err)
		}
		return f
	}
	return nil
}

func StopCPUProfile(f *os.File) {
	if DebugProfile {
		if f != nil {
			defer f.Close()
		}
		defer pprof.StopCPUProfile()
	}
}

func DumpMEMProfile(name string) {
	if DebugProfile {
		f, err := cxcore.CXCreateFile(fmt.Sprintf("%s_%s_mem.pprof", os.Args[0], name))
		if err != nil {
			fmt.Println("Failed to create MEM profile: ", err)
		}
		defer f.Close()
		runtime.GC()
		if err := pprof.WriteHeapProfile(f); err != nil {
			fmt.Println("Failed to write MEM profile: ", err)
		}
	}
}

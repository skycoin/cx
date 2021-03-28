// +build os

package cxos

import (
	"fmt"
	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/util"
	"os"
	"runtime"
	"runtime/pprof"
)

var openProfiles map[string]*os.File = make(map[string]*os.File, 0)

func startCPUProfile(name string, rate int) *os.File {
	f, err := util.CXCreateFile(fmt.Sprintf("%s_%s_cpu.pprof", os.Args[0], name))
	if err != nil {
		fmt.Println("Failed to create CPU profile: ", err)
	}
	if rate != 100 { // test against default value to avoid warning
		runtime.SetCPUProfileRate(rate)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		fmt.Println("Failed to start CPU profile: ", err)
	}
	return f
}

func stopCPUProfile(f *os.File) {
	if f != nil {
		defer f.Close()
	}
	defer pprof.StopCPUProfile()
}

func opStartProfile(inputs []ast.CXValue, outputs []ast.CXValue) {
	profilePath := inputs[0].Get_str()
	openProfiles[profilePath] = startCPUProfile(profilePath, int(inputs[1].Get_i32()))
}

func opStopProfile(inputs []ast.CXValue, outputs []ast.CXValue) {
	profilePath := inputs[0].Get_str()
	stopCPUProfile(openProfiles[profilePath])
}

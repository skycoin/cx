// +build base

package cxcore

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"

	. "github.com/skycoin/cx/cx"
)

var openProfiles map[string]*os.File = make(map[string]*os.File, 0)

func startCPUProfile(name string, rate int) *os.File {
	f, err := CXCreateFile(fmt.Sprintf("%s_%s_cpu.pprof", os.Args[0], name))
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

func opStartProfile(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	profilePath := ReadStr(fp, expr.Inputs[0])
	openProfiles[profilePath] = startCPUProfile(profilePath, int(ReadI32(fp, expr.Inputs[1])))
}

func opStopProfile(prgrm *CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	profilePath := ReadStr(fp, expr.Inputs[0])
	stopCPUProfile(openProfiles[profilePath])
}

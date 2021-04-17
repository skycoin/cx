package cxprof

import (
	"fmt"
	"github.com/skycoin/cx/cx/util"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/sirupsen/logrus"
)

type StopFunc func() time.Duration

func StartProfile(log logrus.FieldLogger) (start time.Time, stop StopFunc) {
	start = time.Now()
	if log != nil {
		log = log.WithField("start", start)
		log.Debug("Started profile.")
	}

	stop = func() time.Duration {
		delta := time.Since(start)
		if log != nil {
			log.WithField("elapsed", delta).Debug("Stopped profile.")
		}
		return delta
	}

	return start, stop
}

func StartCPUProfile(name string, profRate int) (func() error, error) {
	if profRate == 0 {
		return func() error { return nil }, nil
	}

	f, err := util.CXCreateFile(fmt.Sprintf("%s_%s_cpu.pprof", os.Args[0], name))
	if err != nil {
		return func() error { return nil }, err
	}

	// test against default value to avoid warning
	if profRate != 100 {
		runtime.SetCPUProfileRate(profRate)
	}

	if err := pprof.StartCPUProfile(f); err != nil {
		return func() error { return nil }, err
	}

	stop := func() error {
		pprof.StopCPUProfile()
		return f.Close()
	}

	return stop, nil
}

func DumpMemProfile(name string) (err error) {
	var f *os.File
	if f, err = util.CXCreateFile(fmt.Sprintf("%s_%s_mem.pprof", os.Args[0], name)); err != nil {
		return fmt.Errorf("failed to create MEM pprof file: %w", err)
	}

	defer func() {
		if closeErr := f.Close(); err == nil {
			err = closeErr
		}
	}()

	runtime.GC()
	if err = pprof.WriteHeapProfile(f); err != nil {
		return fmt.Errorf("failed to write MEM profile: %w", err)
	}

	return err
}

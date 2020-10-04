package cxflags

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"

	cxcore "github.com/SkycoinProject/cx/cx"
)

// MemoryFlags contains cli flags associated with cx memory parameters.
type MemoryFlags struct {
	initHeapSize string
	maxHeapSize  string
	stackSize    string

	minHeapFreeRatio float64
	maxHeapFreeRatio float64
}

// DefaultMemoryFlags returns the default set of memory flags.
func DefaultMemoryFlags() *MemoryFlags {
	return &MemoryFlags{
		initHeapSize:     strconv.Itoa(cxcore.INIT_HEAP_SIZE),
		maxHeapSize:      strconv.Itoa(cxcore.MAX_HEAP_SIZE),
		stackSize:        strconv.Itoa(cxcore.STACK_SIZE),
		minHeapFreeRatio: float64(cxcore.MIN_HEAP_FREE_RATIO),
		maxHeapFreeRatio: float64(cxcore.MAX_HEAP_FREE_RATIO),
	}
}

// Register registers the flags to a given flag set.
func (mf *MemoryFlags) Register(fs *flag.FlagSet) {
	fs.StringVar(&mf.initHeapSize, "heap-initial", mf.initHeapSize, "initial heap `BYTES` for CX virtual machine (format: <bytes>[G|M|K] )")
	fs.StringVar(&mf.initHeapSize, "hi", mf.initHeapSize, "shorthand for 'heap-size'")
	fs.StringVar(&mf.maxHeapSize, "heap-max", mf.maxHeapSize, "max heap `BYTES` for CX virtual machine (format: <bytes>[G|M|K] )")
	fs.StringVar(&mf.maxHeapSize, "hm", mf.maxHeapSize, "shorthand for 'heap-max'")
	fs.StringVar(&mf.stackSize, "stack-size", mf.stackSize, "stack size in `BYTES` for CX virtual machine (format: <bytes>[G|M|K] )")
	fs.StringVar(&mf.stackSize, "ss", mf.stackSize, "shorthand for 'stack-size'")
	fs.Float64Var(&mf.minHeapFreeRatio, "min-heap-free", mf.minHeapFreeRatio, "`DECIMAL` ratio of the min free heap size after calling GC (range: 0.0 - 1.0)")
	fs.Float64Var(&mf.maxHeapFreeRatio, "max-heap-free", mf.maxHeapFreeRatio, "`DECIMAL` ratio of the max free heap size after calling GC (range: 0.0 - 1.0)")
}

// PostProcess should be called after flags are parsed.
func (mf *MemoryFlags) PostProcess() error {
	var err error

	// Initial heap size.
	cxcore.INIT_HEAP_SIZE, err = parseMemoryString(mf.initHeapSize)
	if err != nil {
		return fmt.Errorf("failed to parse flag 'heap-initial': %w", err)
	}

	// Max heap size.
	cxcore.MAX_HEAP_SIZE, err = parseMemoryString(mf.maxHeapSize)
	if err != nil {
		return fmt.Errorf("failed to parse flag 'heap-max': %w", err)
	}
	if cxcore.INIT_HEAP_SIZE > cxcore.MAX_HEAP_SIZE {
		cxcore.INIT_HEAP_SIZE = cxcore.MAX_HEAP_SIZE
	}

	// Stack size.
	cxcore.STACK_SIZE, err = parseMemoryString(mf.stackSize)
	if err != nil {
		return fmt.Errorf("failed to parse flag 'stack-size': %w", err)
	}

	// Min heap free ratio.
	if err := checkRatio(mf.minHeapFreeRatio); err != nil {
		return fmt.Errorf("failed to parse flag 'min-heap-free': %w", err)
	}
	cxcore.MIN_HEAP_FREE_RATIO = float32(mf.minHeapFreeRatio)

	// Max heap free ratio.
	if err := checkRatio(mf.maxHeapFreeRatio); err != nil {
		return fmt.Errorf("failed to parse flag 'max-heap-free': %w", err)
	}
	cxcore.MAX_HEAP_FREE_RATIO = float32(mf.maxHeapFreeRatio)

	return nil
}

func checkRatio(r float64) error {
	if r < 0 {
		return fmt.Errorf("ratio cannot be smaller than 0.0 (%v)", r)
	}
	if r > 1 {
		return fmt.Errorf("ratio cannot be greater than 1.0 (%v)", r)
	}
	return nil
}

// parseMemoryString is used for the -heap-initial, -heap-max and -stack-size flags.
// This function parses, for example, "1M" to 1048576 (the corresponding number of bytes)
// Possible suffixes are: G or g (gigabytes), M or m (megabytes), K or k (kilobytes)
// Input 'v' is set to 'n' on return if no error occurs.
func parseMemoryString(s string) (int, error) {
	s = strings.TrimSpace(s)

	switch len(s) {
	case 0:
		return 0, errors.New("value cannot be empty")
	case 1:
		n, err := strconv.ParseInt(s, 10, 64)
		return int(n), err
	default:
	}

	switch num, suffix := s[:len(s)-1], s[len(s)-1]; suffix {
	case 'G', 'g':
		return applyFactor(num, 1073741824)
	case 'M', 'm':
		return applyFactor(num, 1048576)
	case 'K', 'k':
		return applyFactor(num, 1024)
	default:
		return applyFactor(num, 1)
	}
}

func applyFactor(num string, fac int) (int, error) {
	n, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return 0, err
	}
	return int(n) * fac, nil
}
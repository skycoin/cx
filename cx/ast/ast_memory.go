package ast

import (
	"fmt"

	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
)

// MinHeapSize determines what's the minimum heap size that a CX program
// needs to have based on INIT_HEAP_SIZE, MAX_HEAP_SIZE and NULL_HEAP_ADDRESS_OFFSET.
func MinHeapSize() types.Pointer {
	minHeapSize := constants.INIT_HEAP_SIZE
	if constants.MAX_HEAP_SIZE < constants.INIT_HEAP_SIZE {
		// Then MAX_HEAP_SIZE overrides INIT_HEAP_SIZE's value.
		minHeapSize = constants.MAX_HEAP_SIZE
	}
	if minHeapSize < constants.NULL_HEAP_ADDRESS_OFFSET {
		// Then the user is trying to allocate too little heap memory.
		// We need at least NULL_HEAP_ADDRESS_OFFSET bytes for `nil`.
		minHeapSize = constants.NULL_HEAP_ADDRESS_OFFSET
	}

	return minHeapSize
}

// EnsureHeap ensures that `prgrm` has `minHeapSize()`
// bytes allocated after the data segment.
func (cxprogram *CXProgram) EnsureMinimumHeapSize() {
	currHeapSize := types.Cast_int_to_ptr(len(cxprogram.Memory)) - cxprogram.Heap.StartsAt
	minHeapSize := MinHeapSize()
	if currHeapSize < minHeapSize {
		cxprogram.Memory = append(cxprogram.Memory, make([]byte, minHeapSize-currHeapSize)...)
		cxprogram.Heap.Size = minHeapSize
	}
}

// ResizeMemory ...
func ResizeMemory(prgrm *CXProgram, newMemSize types.Pointer, isExpand bool) {
	// We can't expand memory to a value greater than `memLimit`.
	if newMemSize > constants.MAX_HEAP_SIZE {
		newMemSize = constants.MAX_HEAP_SIZE
	}

	if newMemSize == prgrm.Heap.Size {
		// Then we're at the limit; we can't expand anymore.
		// We can only hope that the free memory is enough for the CX program to continue running.
		return
	}

	if isExpand {
		// Adding bytes to reach a heap equal to `newMemSize`.
		prgrm.Memory = append(prgrm.Memory, make([]byte, newMemSize-prgrm.Heap.Size)...)
		prgrm.Heap.Size = newMemSize
	} else {
		// Removing bytes to reach a heap equal to `newMemSize`.
		prgrm.Memory = append([]byte(nil), prgrm.Memory[:prgrm.Heap.StartsAt+newMemSize]...)
		prgrm.Heap.Size = newMemSize
	}
}

// AllocateSeq allocates memory in the heap
func AllocateSeq(program interface{}, size types.Pointer) (offset types.Pointer) {
	prgrm, ok := program.(*CXProgram)
	if !ok {
		panic(fmt.Sprintf("error getting cx program"))
	}
	// Current object trying to be allocated would use this address.
	addr := prgrm.Heap.Pointer
	// Next object to be allocated will use this address.
	newFree := addr + size

	// Checking if we can allocate the entirety of the object in the current heap.
	if newFree > prgrm.Heap.Size {
		// It does not fit, so calling garbage collector.
		MarkAndCompact(prgrm)
		// Heap pointer got moved by GC and recalculate these variables based on the new pointer.
		addr = prgrm.Heap.Pointer
		newFree = addr + size

		// If the new heap pointer exceeds `MAX_HEAP_SIZE`, there's nothing left to do.
		if newFree > constants.MAX_HEAP_SIZE {
			panic(constants.HEAP_EXHAUSTED_ERROR)
		}

		// According to MIN_HEAP_FREE_RATIO and MAX_HEAP_FREE_RATION we can either shrink
		// or expand the heap to maintain "healthy" heap sizes. The idea is that we don't want
		// to have an absurdly amount of free heap memory, as we would be wasting resources, and we
		// don't want to have a small amount of heap memory left as we'd be calling the garbage collector
		// too frequently.

		// Calculating free heap memory percentage.
		usedPerc := float32(newFree) / float32(prgrm.Heap.Size)
		freeMemPerc := 1.0 - usedPerc

		// Then we have less than MIN_HEAP_FREE_RATIO memory left. Expand!
		if freeMemPerc < constants.MIN_HEAP_FREE_RATIO {
			// Calculating new heap size in order to reach MIN_HEAP_FREE_RATIO.
			newMemSize := types.Cast_f32_to_ptr(float32(newFree) / (1.0 - constants.MIN_HEAP_FREE_RATIO))
			ResizeMemory(prgrm, newMemSize, true)
		}

		// Then we have more than MAX_HEAP_FREE_RATIO memory left. Shrink!
		if freeMemPerc > constants.MAX_HEAP_FREE_RATIO {
			// Calculating new heap size in order to reach MAX_HEAP_FREE_RATIO.
			newMemSize := types.Cast_f32_to_ptr(float32(newFree) / (1.0 - constants.MAX_HEAP_FREE_RATIO))

			// This check guarantees that the CX program has always at least INIT_HEAP_SIZE bytes to work with.
			// A flag could be added later to remove this, as in some cases this mechanism could not be desired.
			if newMemSize > constants.INIT_HEAP_SIZE {
				ResizeMemory(prgrm, newMemSize, false)
			}
		}
	}

	prgrm.Heap.Pointer = newFree

	// Returning absolute memory address (not relative to where heap starts at).
	// Above this point we were performing all operations taking into
	// consideration only heap offsets.
	return addr + prgrm.Heap.StartsAt
}

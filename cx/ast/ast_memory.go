package ast

import (
    "github.com/skycoin/cx/cx/constants"
    "github.com/skycoin/cx/cx/types"
)

// minHeapSize determines what's the minimum heap size that a CX program
// needs to have based on INIT_HEAP_SIZE, MAX_HEAP_SIZE and NULL_HEAP_ADDRESS_OFFSET.
func minHeapSize() types.Pointer {
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
	currHeapSize := types.Cast_int_to_ptr(len(cxprogram.Memory)) - cxprogram.HeapStartsAt
	minHeapSize := minHeapSize()
	if currHeapSize < minHeapSize {
		cxprogram.Memory = append(cxprogram.Memory, make([]byte, minHeapSize-currHeapSize)...)
	}
}


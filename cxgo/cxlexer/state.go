package cxlexer

import (
	cxcore "github.com/SkycoinProject/cx/cx"
	"github.com/SkycoinProject/cx/cxgo/actions"
)

// ChainStatePrelude initializes the program structure `prog` with data from
// the program state stored on a CX chain.
// This is used for transaction/broadcast mode.
// TODO @evanlinjin: Fill this out.
func ChainStatePrelude(prog *cxcore.CXProgram) (state, heap []byte) {
	return nil, nil
}

// MergeBlockchainHeap adds the heap `bcHeap` found in the program state of a CX
// chain to the program to be run `PRGRM` and updates all the references to heap
// objects found in the transaction code considering the data segment found in
// the serialized program `sPrgrm`.
func MergeBlockchainHeap(bcHeap, sPrgrm []byte) {
	// Setting the CX runtime to run `PRGRM`.
	actions.PRGRM.SelectProgram()

	bcHeapLen := len(bcHeap)
	remHeapSpace := len(actions.PRGRM.Memory[actions.PRGRM.HeapStartsAt:])
	fullDataSegSize := actions.PRGRM.HeapStartsAt - actions.PRGRM.StackSize
	// Copying blockchain code heap.
	if bcHeapLen > remHeapSpace {
		// We don't have enough space. We're using the available bytes...
		for c := 0; c < remHeapSpace; c++ {
			actions.PRGRM.Memory[actions.PRGRM.HeapStartsAt+c] = bcHeap[c]
		}
		// ...and then we append the remaining bytes.
		actions.PRGRM.Memory = append(actions.PRGRM.Memory, bcHeap[remHeapSpace:]...)
	} else {
		// We have enough space and we simply write the bytes.
		for c := 0; c < bcHeapLen; c++ {
			actions.PRGRM.Memory[actions.PRGRM.HeapStartsAt+c] = bcHeap[c]
		}
	}
	// Recalculating the heap size.
	actions.PRGRM.HeapSize = len(actions.PRGRM.Memory) - actions.PRGRM.HeapStartsAt
	txnDataLen := fullDataSegSize - cxcore.GetSerializedDataSize(sPrgrm)
	// TODO: CX chains only work with one package at the moment (in the blockchain code). That is what that "1" is for.
	// Displacing the references to heap objects by `txnDataLen`.
	// This needs to be done as the addresses to the heap objects are displaced
	// by the addition of the transaction code's data segment.
	cxcore.DisplaceReferences(actions.PRGRM, txnDataLen, 1)
}
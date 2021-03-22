package cxlexer

import (
	"github.com/skycoin/cx/cx"
	"github.com/skycoin/cx/cxgo/actions"
)

// InitProg initiates program contained in 'actions.PRGRM'.
func InitProg() (*cxcore.CXProgram, error) {
	coreProgState, err := cxcore.GetProgram()
	if err != nil {
		return nil, err
	}
	prog := cxcore.MakeProgram()
	prog.Packages = coreProgState.Packages
	actions.PRGRM = prog
	return prog, nil
}

// ProgBytes represents program bytes.
type ProgBytes struct {
	State []byte
	Heap  []byte
}

// LoadProgFromBytes initializes the program structure `prog` with data from
// the program state stored on a CX chain.
// This is used for transaction/broadcast mode.
func LoadProgFromBytes(prog *cxcore.CXProgram, progS []byte) (*ProgBytes, error) {
	memOffset := cxcore.GetSerializedMemoryOffset(progS)
	stackSize := cxcore.GetSerializedStackSize(progS)

	// program state with only memory stack and heap
	progSMem := progS[:memOffset]

	// append new stack
	progSMem = append(progSMem, make([]byte, stackSize)...)

	// append data and heap segment
	progSMem = append(progSMem, progS[memOffset:]...)

	heap := progS[memOffset+cxcore.GetSerializedDataSize(progS):]

	*prog = *cxcore.Deserialize(progSMem)
	actions.PRGRM = prog
	actions.DataOffset = prog.HeapStartsAt // Start adding data elements here.

	return &ProgBytes{
		State: progS,
		Heap:  heap,
	}, nil
}

/*
// MergeChainHeap adds the heap `bcHeap` found in the program state of a CX
// chain to the program to be run `PRGRM` and updates all the references to heap
// objects found in the transaction code considering the data segment found in
// the serialized program `sPrgrm`.
func (pb *ProgBytes) MergeChainHeap() error {
	// Setting the CX runtime to run `PRGRM`.
	if _, err := actions.PRGRM.SelectProgram(); err != nil {
		return err
	}

	bcHeapLen := len(pb.Heap)
	remHeapSpace := len(actions.PRGRM.Memory[actions.PRGRM.HeapStartsAt:])
	fullDataSegSize := actions.PRGRM.HeapStartsAt - actions.PRGRM.StackSize
	// Copying blockchain code heap.
	if bcHeapLen > remHeapSpace {
		// We don't have enough space. We're using the available bytes...
		for c := 0; c < remHeapSpace; c++ {
			actions.PRGRM.Memory[actions.PRGRM.HeapStartsAt+c] = pb.Heap[c]
		}
		// ...and then we append the remaining bytes.
		actions.PRGRM.Memory = append(actions.PRGRM.Memory, pb.Heap[remHeapSpace:]...)
	} else {
		// We have enough space and we simply write the bytes.
		for c := 0; c < bcHeapLen; c++ {
			actions.PRGRM.Memory[actions.PRGRM.HeapStartsAt+c] = pb.Heap[c]
		}
	}
	// Recalculating the heap size.
	actions.PRGRM.HeapSize = len(actions.PRGRM.Memory) - actions.PRGRM.HeapStartsAt
	txnDataLen := fullDataSegSize - cxcore.GetSerializedDataSize(pb.State)
	// TODO: CX chains only work with one package at the moment (in the blockchain code). That is what that "1" is for.
	// Displacing the references to heap objects by `txnDataLen`.
	// This needs to be done as the addresses to the heap objects are displaced
	// by the addition of the transaction code's data segment.
	cxcore.DisplaceReferences(actions.PRGRM, txnDataLen, 1)

	return nil
}
*/
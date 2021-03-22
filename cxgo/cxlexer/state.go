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

/*
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
*/
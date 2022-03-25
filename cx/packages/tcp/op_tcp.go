package tcp

import (
	"context"
	"log"
	"net"
	"net/rpc"
	"time"

	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/types"
)

//todo need to find a way to support Listener and connection interface
var ln net.Listener

var lc net.ListenConfig

var conn net.Conn

var DefaultServer = rpc.NewServer()

func opTCPDial(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	network, address, errorstring := inputs[0].Arg, inputs[1].Arg, outputs[0].Arg
	fp := inputs[0].FramePointer

	log.Println("network", network)
	log.Println("address", address)

	var err error
	conn, err = net.Dial("tcp", "localhost:9000")

	if err != nil {
		types.Write_str(prgrm, prgrm.Memory, ast.GetFinalOffset(prgrm, fp, errorstring), err.Error())
	}

	conn.Close()
}

func opTCPClose(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	ln.Close()
}

func opTCPAccept(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	conn, _ = ln.Accept()
	conn.Close()
}

func opTCPListen(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
	network, address, errorstring := inputs[0].Arg, inputs[1].Arg, outputs[0].Arg
	fp := inputs[0].FramePointer

	log.Println("network", network)
	log.Println("address", address)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)

	defer cancel()

	var err error

	ln, err = lc.Listen(ctx, "tcp", ":9000")

	ln.Close()

	if err != nil {
		types.Write_str(prgrm, prgrm.Memory, ast.GetFinalOffset(prgrm, fp, errorstring), err.Error())
	}
}

package opcodes

import (
	"context"
	"github.com/skycoin/cx/cx/ast"
	"log"
	"net"
	"net/rpc"
	"time"
)

//todo need to find a way to support Listener and connection interface
var ln net.Listener

var lc net.ListenConfig

var conn net.Conn

var DefaultServer = rpc.NewServer()

func init() {
	netPkg := ast.MakePackage("tcp")

	dialerStrct := ast.MakeStruct("Dialer")

	netPkg.AddStruct(dialerStrct)

	ast.PROGRAM.AddPackage(netPkg)
}

func opTCPDial(inputs []ast.CXValue, outputs []ast.CXValue) {
	network, address, errorstring := inputs[0].Arg, inputs[1].Arg, outputs[0].Arg
    fp := inputs[0].FramePointer

	log.Println("network", network)
	log.Println("address", address)

	var err error
	conn, err = net.Dial("tcp", "localhost:9000")

	if err != nil {
		ast.WriteString(fp, err.Error(), errorstring)
	}

	conn.Close()
}

func opTCPClose(inputs []ast.CXValue, outputs []ast.CXValue) {
	ln.Close()
}

func opTCPAccept(inputs []ast.CXValue, outputs []ast.CXValue) {
	conn, _ = ln.Accept()
	conn.Close()
}

func opTCPListen(inputs []ast.CXValue, outputs []ast.CXValue) {
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
		ast.WriteString(fp, err.Error(), errorstring)
	}
}

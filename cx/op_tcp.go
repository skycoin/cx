package cxcore

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

func opTCPDial(expr *ast.CXExpression, fp int) {

	network, address, errorstring := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]

	log.Println("network", network)
	log.Println("address", address)

	conn, err := net.Dial("tcp", "localhost:9000")

	if err != nil {
		ast.WriteString(fp, err.Error(), errorstring)
	}

	conn.Close()

}

func opTCPClose(expr *ast.CXExpression, fp int) {

	ln.Close()
}

func opTCPAccept(expr *ast.CXExpression, fp int) {

	conn, _ = ln.Accept()

	conn.Close()
}

func opTCPListen(expr *ast.CXExpression, fp int) {

	network, address, errorstring := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]

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

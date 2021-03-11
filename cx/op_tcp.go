package cxcore

import (
	"context"
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

	netPkg := MakePackage("tcp")

	dialerStrct := MakeStruct("Dialer")

	netPkg.AddStruct(dialerStrct)

	PROGRAM.AddPackage(netPkg)

}

func opTCPDial(expr *CXExpression, fp int) {

	network, address, errorstring := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]

	log.Println("network", network)
	log.Println("address", address)

	conn, err := net.Dial("tcp", "localhost:9000")

	if err != nil {
		WriteString(fp, err.Error(), errorstring)
	}

	conn.Close()

}

func opTCPClose(expr *CXExpression, fp int) {

	ln.Close()
}

func opTCPAccept(expr *CXExpression, fp int) {

	conn, _ = ln.Accept()

	conn.Close()
}

func opTCPListen(expr *CXExpression, fp int) {

	network, address, errorstring := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]

	log.Println("network", network)
	log.Println("address", address)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)

	defer cancel()

	var err error

	ln, err = lc.Listen(ctx, "tcp", ":9000")

	ln.Close()

	if err != nil {
		WriteString(fp, err.Error(), errorstring)
	}

}

package cxcore

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"time"
)

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

	log.Println("opTCPDial")
	network, address, errorstring := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]

	log.Println("opTCPDial")

	fmt.Println("network", network)
	fmt.Println("address", address)

	//	url := ReadStr(fp, network)

	conn, err := net.Dial("tcp", "localhost:9000")

	if err != nil {
		WriteString(fp, err.Error(), errorstring)
	}

	conn.Close()

}

//

func opTCPClose(expr *CXExpression, fp int) {

	ln.Close()
}

func opTCPAccept(expr *CXExpression, fp int) {

	var err error

	conn, err = ln.Accept()

	fmt.Println("ln", conn)

	fmt.Println("err", err)
}

func opTCPListen(expr *CXExpression, fp int) {

	network, address, errorstring := expr.Inputs[0], expr.Inputs[1], expr.Outputs[0]

	fmt.Println("network", network)
	fmt.Println("address", address)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)

	defer cancel()

	var err error

	ln, err = lc.Listen(ctx, "tcp", ":9000")

	ln.Close()

	fmt.Println("ln", ln)

	fmt.Println("network", network)
	fmt.Println("address", address)

	fmt.Println("errorstring", errorstring)

	if err != nil {
		WriteString(fp, err.Error(), errorstring)
	}

}

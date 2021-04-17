package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/skycoin/cx/cmd/cxplayground/playground"
	"github.com/skycoin/cx/cxparser/actions"
	"github.com/skycoin/cx/cxparser/webapi"
)

func main() {
	host := ":5336"

	mux := http.NewServeMux()
	workingDir, _ := os.Getwd()
	if err := playground.InitPlayground(workingDir); err != nil {
		// error captured while initiating the playground examples, should be handled in the future
		fmt.Println("Fail to initiating palyground examples")
	}

	mux.HandleFunc("/playground/examples", playground.GetExampleFileList)
	mux.HandleFunc("/playground/examples/code", playground.GetExampleFileContent)

	mux.Handle("/", http.FileServer(http.Dir("./dist")))
	mux.Handle("/program/", webapi.NewAPI("/program", actions.AST))
	mux.HandleFunc("/eval", playground.RunProgram)

	if listener, err := net.Listen("tcp", host); err == nil {
		fmt.Println("Starting web service for CX playground on http://127.0.0.1:5336/")
		http.Serve(listener, mux)
	}
}

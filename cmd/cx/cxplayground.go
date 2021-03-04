package main

import (
	"encoding/json"
	"fmt"
	"github.com/skycoin/cx/cxgo/actions"
	"github.com/skycoin/cx/cxgo/api"
	"io/ioutil"
	"net"
	"net/http"
)

func cxPlayground() {
	host := ":5336"

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("./dist")))
	mux.Handle("/program/", api.NewAPI("/program", actions.PRGRM))
	mux.HandleFunc("/eval", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var b []byte
		var err error
		if b, err = ioutil.ReadAll(r.Body); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		var source SourceCode
		if err := json.Unmarshal(b, &source); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		if err := r.ParseForm(); err == nil {
			fmt.Fprintf(w, "%s", Eval(source.Code+"\n"))
		}
	})

	if listener, err := net.Listen("tcp", host); err == nil {
		fmt.Println("Starting web service for CX playground on http://127.0.0.1:5336/")
		http.Serve(listener, mux)
	}
}

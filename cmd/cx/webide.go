package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/skycoin/cx/cxgo/actions"
	"github.com/skycoin/cx/cxgo/api"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

//web interactive mode

func ServiceMode() {
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
		fmt.Println("Starting CX web service on http://127.0.0.1:5336/")
		http.Serve(listener, mux)
	}
}

func WebIdeServiceMode() {
	// Leaps's host address
	ideHost := "localhost:5335"

	// Working directory for ide
	sharedPath := fmt.Sprintf("%s/src/github.com/skycoin/cx", os.Getenv("GOPATH"))

	// Start Leaps
	// cmd = `leaps -address localhost:5335 $GOPATH/src/skycoin/cx`
	cmnd := exec.Command("leaps", "-address", ideHost, sharedPath)

	// Just leave start command
	cmnd.Start()
}

func PersistentServiceMode() {
	fmt.Println("Start persistent for service mode!")

	fi := bufio.NewReader(os.Stdin)

	for {
		var inp string
		var ok bool

		printPrompt()

		if inp, ok = readline(fi); ok {
			if isJSON(inp) {
				var err error
				client := &http.Client{}
				body := bytes.NewBufferString(inp)
				req, err := http.NewRequest("GET", "http://127.0.0.1:5336/eval", body)
				if err != nil {
					fmt.Println(err)
					return
				}

				if resp, err := client.Do(req); err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(resp.Status)
				}
			}
		}
	}
}

func isJSON(str string) bool {
	var js map[string]interface{}
	err := json.Unmarshal([]byte(str), &js)
	return err == nil
}

// ----------------------------------------------------------------
//                     Utility functions

func readline(fi *bufio.Reader) (string, bool) {
	s, err := fi.ReadString('\n')

	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, "\r", "", -1)

	for _, ch := range s {
		if ch == rune(4) {
			err = io.EOF
			break
		}
	}

	if err != nil {
		return "", false
	}

	return s, true
}

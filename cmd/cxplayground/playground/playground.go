package playground

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/execute"

	"github.com/skycoin/cx/cxparser/actions"
	cxparsing "github.com/skycoin/cx/cxparser/cxparsing"

	cxparsingcompletor "github.com/skycoin/cx/cxparser/cxparsingcompletor"
	cxpartialparsing "github.com/skycoin/cx/cxparser/cxpartialparsing"
)

var (
	exampleCollection map[string]string
	exampleNames      []string
	examplesDir       string
)

type ExampleContent struct {
	ExampleName string
}

var InitPlayground = func(workingDir string) error {
	examplesDir = filepath.Join(workingDir, "../../examples")
	exampleCollection = make(map[string]string)

	exampleInfoList, err := ioutil.ReadDir(examplesDir)
	if err != nil {
		fmt.Printf("Fail to get file list under examples dir: %s\n", err.Error())
		return err
	}
	for _, exp := range exampleInfoList {
		if exp.IsDir() {
			continue
		}
		path := filepath.Join(examplesDir, exp.Name())
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Printf("Fail to read example file %s\n", path)
			// coninue if fail to read the current example file
			continue
		}

		exampleName := exp.Name()
		exampleNames = append(exampleNames, exampleName)
		exampleCollection[exampleName] = string(bytes)
	}

	return nil
}

func GetExampleFileList(w http.ResponseWriter, r *http.Request) {

	exampleNamesBytes, err := json.Marshal(exampleNames)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, string(exampleNamesBytes))
}

func GetExampleFileContent(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		b   []byte
		err error
	)
	if b, err = ioutil.ReadAll(r.Body); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var example ExampleContent
	if err = json.Unmarshal(b, &example); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	content, err := getExampleFileContent(example.ExampleName)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, content)
}

var getExampleFileContent = func(exampleName string) (string, error) {
	exampleContent, ok := exampleCollection[exampleName]
	if !ok {
		err := fmt.Errorf("example name %s not found", exampleName)

		return "", err
	}
	return exampleContent, nil
}

type SourceCode struct {
	Code string `json:"code,omitempty"`
}

func RunProgram(w http.ResponseWriter, r *http.Request) {
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
		fmt.Fprintf(w, "%s", eval(source.Code+"\n"))
	}
}

func unsafeeval(code string) (out string) {
	var lexer *cxparsingcompletor.Lexer
	defer func(lexer *cxparsingcompletor.Lexer) {
		if r := recover(); r != nil {
			out = fmt.Sprintf("%v", r)
			// lexer.Stop()
			return
		}
	}(lexer)

	// storing strings sent to standard output
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		return fmt.Sprintf("%v", err)
	}
	os.Stdout = w

	actions.LineNo = 0

	actions.AST = ast.MakeProgram()
	cxpartialparsing.Program = actions.AST

	cxpartialparsing.Parse(code)

	actions.AST = cxpartialparsing.Program

	lexer = cxparsingcompletor.NewLexer(bytes.NewBufferString(code))
	cxparsingcompletor.Parse(lexer)
	//yyParse(lexer)

	err = cxparsing.AddInitFunction(actions.AST)
	if err != nil {
		return fmt.Sprintf("%s", err)
	}
	//if err := actions.AST.RunCompiled(0, nil); err != nil {
	err = execute.RunCompiled(actions.AST, 0, nil)
	if err != nil {
		actions.AST = ast.MakeProgram()
		return fmt.Sprintf("%s", err)
	}

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	w.Close()
	os.Stdout = old // restoring the real stdout
	out = <-outC

	actions.AST = ast.MakeProgram()
	return out
}

func eval(code string) string {
	runtime.GOMAXPROCS(2)
	ch := make(chan string, 1)

	var result string

	go func() {
		result = unsafeeval(code)
		ch <- result
	}()

	timer := time.NewTimer(20 * time.Second)
	defer timer.Stop()

	select {
	case <-ch:
		return result
	case <-timer.C:
		actions.AST = ast.MakeProgram()
		return "Timed out."
	}
}

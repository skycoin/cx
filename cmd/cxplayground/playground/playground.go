package playground

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

var exampleCollection map[string]string
var exampleNames []string
var examplesDir string

type ExampleContent struct {
	ExampleName string
}

func InitPlayground(workingDir string) error {
	examplesDir = filepath.Join(workingDir, "examples")
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
	var b []byte
	var err error
	if b, err = ioutil.ReadAll(r.Body); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Println(string(b))
	var example ExampleContent
	if err := json.Unmarshal(b, &example); err != nil {
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

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

func GetPlayground(w http.ResponseWriter, r *http.Request) {
	playgroundHtml := `<html>

<head>
    <title>test</title>
    <script src="https://libs.baidu.com/jquery/1.7.2/jquery.min.js"></script>
</head>
<script>
    // var url = "http://localhost:5336/"
    $().ready(function () {
        $.getJSON("playground/examples", function (inputData) {
            $.each(inputData, function (i) {
                $("#list").append("<option value='" + i + "'>" + inputData[i] + "</option>");
            });
        });
        $("#list").bind("change", function () {
            var data = { "examplename": $("#list").find("option:selected").text() };
            $.ajax({
                type: "POST",
                url: "/playground/examples/code",
                contentType: "application/json;charset=utf-8",
                data: JSON.stringify(data),
                cache: false,
                success: function (message) {
                    $("#inputData").html(message);

                },
                error: function (message) {
                    $("#inputData").html(message);
                }
            });
        });
        $("#run").click(function () {
            var data = { "code": $("#inputData").text() };
            $.ajax({
                type: "POST",
                url: "/eval",
                contentType: "application/json;charset=utf-8",
                data: JSON.stringify(data),
                cache: false,
                success: function (message) {
                    $("#result").text(function (i, origText) {
                        return message;
                    });
                },
                error: function (message) {
                    $("#result").text(function (i, origText) {
                        return message;
                    });
                }
            });
        });
    });

</script>

<body>
    <div style="text-align:center;">
        <div id="banner">
            <div id="head" itemprop="name">Playground</div>
            <select id="list">
                <option value="hello">Hello, World</option>
            </select>
            <input type="button" value="Run" id="run">
        </div>
        <div>
            <textarea id="inputData" rows="20" cols="80">
package main;

func main(){
    str.print("Hello, world!")
}
</textarea>
            <p>result:</p>
            <p id="result"></p>
        </div>
    </div>
</body>

</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fmt.Fprintf(w, playgroundHtml)

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

func getExampleFileContent(exampleName string) (string, error) {
	exampleContent, ok := exampleCollection[exampleName]
	if !ok {
		err := fmt.Errorf("example name %s not found", exampleName)

		return "", err
	}
	return exampleContent, nil
}

package stage1

import (
	"os"
	"path/filepath"
	"testing"
)

var files = map[string]string{
	"main.cx": `
package main
import "package1"
`,
	"package1/test2.cx": `
package package1

import package2
import package3
`,
	"package2/test3.cx": `
package package2
`,
	"package3/test4.cx": `
package package3
`,
	"package3/test5.cx": `
package package3
var a i32
a = 2
`,
}

func setup(testPath string) {
	// create files here
	for file, value := range files {
		file = filepath.Join(testPath, file)
		if f, err := CreateFile(file); err != nil {
			panic(err)
		} else {
			f.WriteString(value)
			f.Close()
		}
	}
}

func CreateFile(fileName string) (*os.File, error) {
	parentFolder := filepath.Dir(fileName)
	err := os.MkdirAll(parentFolder, 0755)
	if err != nil {
		return nil, err
	}
	file, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func shutdown() {
	// delete files here
}

func TestParsePackages(t *testing.T) {
	tests := []struct {
		input    string
		packages []string
	}{
		{"main.cx", []string{"main", "package1", "package2", "package3"}},
	}

	tempDir := t.TempDir()
	setup(tempDir)

	defer shutdown()

	for _, tt := range tests {
		file, err := os.Open(filepath.Join(tempDir, tt.input))
		if err != nil {
			t.Errorf("error: %+v", err)
		}
		packages, err := ParsePackages([]*os.File{file}, []string{filepath.Join(tempDir, tt.input)})
		if err != nil {
			t.Errorf("error parsing packages %+v", err)
		}
		if len(packages) != 4 {
			t.Errorf("expected %d packages, got %d", 4, len(packages))
		}
		for i, pkg := range packages {
			if pkg.Name != tt.packages[i] {
				t.Errorf("error: package mismatch, expected: %s, got: %s", tt.packages[i], pkg.Name)
			}
		}
	}
}

package declaration_extraction

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestDeclarationExtraction(t *testing.T) {

	type Test struct {
		content  string
		fileName string
		GlblDec  []Declaration
		EnumDec  []EnumDeclaration
		StrctDec []Declaration
		FuncDec  []Declaration
	}

	test1 := Test{
		content: `package hello
var apple string

		var banana string

/* sfds*/

//asdsdsa

//sdfsd
type person struct {
	name string
}

func main () {
	var local string

	type animal struct {
		kingdom string
	}
}

	func functionTwo () {

}

type Direction int

const (
	North Direction = iota
	South
	East
	West
)

const (
	First Number = iota
	Second
)`,
		fileName: "test.cx",
		GlblDec: []Declaration{
			{
				PackageID:   "hello",
				FileID:      "test.cx",
				StartOffset: 0,
				Length:      16,
				Name:        "apple",
			},
			{
				PackageID:   "hello",
				FileID:      "test.cx",
				StartOffset: 2,
				Length:      17,
				Name:        "banana",
			},
		},
		EnumDec: []EnumDeclaration{
			{
				PackageID:   "hello",
				FileID:      "test.cx",
				StartOffset: 1,
				Length:      15,
				Name:        "North",
				Type:        "Direction",
				Value:       0,
			},
			{
				PackageID:   "hello",
				FileID:      "test.cx",
				StartOffset: 1,
				Length:      5,
				Name:        "South",
				Type:        "Direction",
				Value:       1,
			},
			{
				PackageID:   "hello",
				FileID:      "test.cx",
				StartOffset: 1,
				Length:      4,
				Name:        "East",
				Type:        "Direction",
				Value:       2,
			},
			{
				PackageID:   "hello",
				FileID:      "test.cx",
				StartOffset: 1,
				Length:      4,
				Name:        "West",
				Type:        "Direction",
				Value:       3,
			},
			{
				PackageID:   "hello",
				FileID:      "test.cx",
				StartOffset: 1,
				Length:      12,
				Name:        "First",
				Type:        "Number",
				Value:       0,
			},
			{
				PackageID:   "hello",
				FileID:      "test.cx",
				StartOffset: 1,
				Length:      6,
				Name:        "Second",
				Type:        "Number",
				Value:       1,
			},
		},
		StrctDec: []Declaration{
			{
				PackageID:   "hello",
				FileID:      "test.cx",
				StartOffset: 0,
				Length:      18,
				Name:        "person",
			},
			{
				PackageID:   "hello",
				FileID:      "test.cx",
				StartOffset: 1,
				Length:      18,
				Name:        "animal",
			},
			{
				PackageID:   "hello",
				FileID:      "test.cx",
				StartOffset: 0,
				Length:      18,
				Name:        "Direction",
			},
		},
		FuncDec: []Declaration{
			{
				PackageID:   "hello",
				FileID:      "test.cx",
				StartOffset: 0,
				Length:      12,
				Name:        "main",
			},
			{
				PackageID:   "hello",
				FileID:      "test.cx",
				StartOffset: 1,
				Length:      19,
				Name:        "functionTwo",
			},
		},
	}

	file, err := os.Create(test1.fileName)
	num, err := file.Write([]byte(test1.content))
	src, err := os.Open(test1.fileName)

	if err != nil {
		log.Fatal(err)
		fmt.Println(num)
	}

	Glbl := extractGlbl(src)

	if Glbl == nil {
		t.Error("No Global Declarations.")
	}

	for i := range Glbl {
		if Glbl[i] != test1.GlblDec[i] {
			t.Errorf("Global Declaration got %+v : want %+v", Glbl[i], test1.GlblDec[i])
		}
	}

	src, err = os.Open(test1.fileName)
	Enum := extractEnum(src)
	if Enum == nil {
		t.Error("No Enum Declarations.")
	}

	for i := range Enum {
		if Enum[i] != test1.EnumDec[i] {
			t.Errorf("Enum Declaration got %+v : want %+v", Enum[i], test1.EnumDec[i])
		}
	}

	src, err = os.Open(test1.fileName)
	Strct := extractStrct(src)
	if Strct == nil {
		t.Error("No Struct Declarations.")
	}

	for i := range Strct {
		if Strct[i] != test1.StrctDec[i] {
			t.Errorf("Struct Declaration got %+v : want %+v", Strct[i], test1.StrctDec[i])
		}
	}

	src, err = os.Open(test1.fileName)
	Func := extractFunc(src)
	if Func == nil {
		t.Error("No Function Declarations.")
	}

	for i := range Func {
		if Func[i] != test1.FuncDec[i] {
			t.Errorf("Function Declaration got %+v : want %+v", Func[i], test1.FuncDec[i])
		}
	}
}

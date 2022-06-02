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
	`,
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
				StartOffset: 8,
				Length:      17,
				Name:        "banana",
			},
		},
		EnumDec: []EnumDeclaration{},
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
				StartOffset: 4,
				Length:      18,
				Name:        "animal",
			},
		},
		FuncDec: []Declaration{
			{
				PackageID:   "hello",
				FileID:      "test.cx",
				StartOffset: 0,
				Length:      9,
				Name:        "main",
			},
			{
				PackageID:   "hello",
				FileID:      "test.cx",
				StartOffset: 0,
				Length:      16,
				Name:        "functionTwo",
			},
		},
	}

	file, err := os.Create(test1.fileName)
	num, err := file.WriteString(test1.content)
	src, err := os.Open(test1.fileName)

	Glbl, err := extractGlbl(src)
	// Enum, err := extractEnum(file)
	// Strct, err := extractStrct(file)
	// Func, err := extractFunc(file)
	if err != nil {
		log.Fatal(err)
		fmt.Println(num)
	}

	for i := range Glbl {
		if Glbl[i] != test1.GlblDec[i] {
			t.Errorf("got %+v : want %+v", Glbl[i], test1.GlblDec[i])
		}

	}
}

func TestExtractStrct(t *testing.T) {

	file, err := os.Open("test.cx")
	got, err := extractStrct(file)

	if err != nil {
		log.Fatal(err)
	}

	want := []Declaration{
		Declaration{
			PackageID:   "hello",
			FileID:      "test.cx",
			StartOffset: 0,
			Length:      18,
			Name:        "person",
		},
		Declaration{
			PackageID:   "hello",
			FileID:      "test.cx",
			StartOffset: 4,
			Length:      18,
			Name:        "animal",
		},
	}

	//Check if any declaration were detected
	if got == nil {
		t.Error("No Struct Declarations")
	}

	for i := range got {
		if got[i] != want[i] {
			t.Errorf("got %+v   want %+v\n", got[i], want[i])
		}
	}

}

func TestExtractFunc(t *testing.T) {

	file, err := os.Open("test.cx")
	got, err := extractFunc(file)

	if err != nil {
		log.Fatal(err)
	}

	want := []Declaration{
		{
			PackageID:   "hello",
			FileID:      "test.cx",
			StartOffset: 0,
			Length:      9,
			Name:        "main",
		},
		{
			PackageID:   "hello",
			FileID:      "test.cx",
			StartOffset: 0,
			Length:      16,
			Name:        "functionTwo",
		},
	}

	//Check if any declaration were detected
	if got == nil {
		t.Error("No Function Declarations")
	}

	for i := range got {
		if got[i] != want[i] {
			t.Errorf("got %+v   want %+v\n", got[i], want[i])
		}

	}

}

func TestRmComment(t *testing.T) {

	file, err := os.ReadFile("test.cx")

	got := rmComment(file)

	if err != nil {
		log.Fatal(err)
	}

	t.Error(string(got))
}

func TestExtractPkg(t *testing.T) {

	file, err := os.ReadFile("./test.cx")

	got := extractPkg(file)

	if err != nil {
		log.Fatal(err)
	}

	t.Error(got)

}

func TestExtractEnum(t *testing.T) {
	file, err := os.Open("test.cx")
	got, err := extractEnum(file)

	if err != nil {
		log.Fatal(err)
	}

	//Check if any declaration were detected
	if got == nil {
		t.Error("No Function Declarations")
	}

	for i := range got {
		t.Errorf(" %+v", got[i])
	}
	// t.Errorf(" %+v", got[0])

}

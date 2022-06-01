package declaration_extraction

import (
	"log"
	"os"
	"testing"
)

func TestExtractGlbl(t *testing.T) {

	got, err := extractGlbl("test.cx")

	if err != nil {
		log.Fatal(err)
	}

	want := []Declaration{
		Declaration{
			PackageID:   "hello",
			FileID:      "test.cx",
			StartOffset: 0,
			Length:      16,
			Name:        "apple",
		},
		Declaration{
			PackageID:   "hello",
			FileID:      "test.cx",
			StartOffset: 8,
			Length:      17,
			Name:        "banana",
		},
	}

	// Check if any declaration were detected
	if got == nil {
		t.Error("No Global Declarations")
	}

	for i := range got {
		if got[i] != want[i] {
			t.Errorf("got %+v   want %+v\n", got[i], want[i])
		}
	}

}

func TestExtractStrct(t *testing.T) {

	got, err := extractStrct("test.cx")

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

	got, err := extractFunc("test.cx")

	if err != nil {
		log.Fatal(err)
	}

	want := []Declaration{
		Declaration{
			PackageID:   "hello",
			FileID:      "test.cx",
			StartOffset: 0,
			Length:      9,
			Name:        "main",
		},
		Declaration{
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

package declaration_extraction

import (
	"log"
	"os"
	"testing"
)

func TestExtractGlbl(t *testing.T) {
	file, err := os.Open("./test.cx")

	if err != nil {
		log.Fatal(err)
	}

	got := extractGlbl(file, "test.cx")

	want := []Declaration{
		Declaration{
			PackageID:   "hello",
			FileID:      "test.cx",
			StartOffset: 0,
			Length:      3,
			Name:        "apple",
		},
		Declaration{
			PackageID:   "hello",
			FileID:      "test.cx",
			StartOffset: 7,
			Length:      3,
			Name:        "banana",
		},
	}

	//Check if any declaration were detected
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
	file, err := os.Open("./test.cx")

	if err != nil {
		log.Fatal(err)
	}

	got := extractStrct(file, "test.cx")

	want := []Declaration{
		Declaration{
			PackageID:   "hello",
			FileID:      "test.cx",
			StartOffset: 0,
			Length:      3,
			Name:        "person",
		},
		Declaration{
			PackageID:   "hello",
			FileID:      "test.cx",
			StartOffset: 3,
			Length:      3,
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
	file, err := os.Open("./test.cx")

	if err != nil {
		log.Fatal(err)
	}

	got := extractFunc(file, "test.cx")

	want := []Declaration{
		Declaration{
			PackageID:   "hello",
			FileID:      "test.cx",
			StartOffset: 0,
			Length:      3,
			Name:        "main",
		},
		Declaration{
			PackageID:   "hello",
			FileID:      "test.cx",
			StartOffset: 0,
			Length:      3,
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

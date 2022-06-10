package encoder

import (
	"testing"

	"github.com/skycoin/cx/cmd/packageloader/bolt"
)

func TestSavePackage(t *testing.T) {
	tests := []struct {
		Scenario string
		Database string
	}{
		{
			"Test with Redis database",
			"redis",
		},
		{
			"Test with Bolt database",
			"bolt",
		},
	}
	for _, testcase := range tests {
		t.Run(testcase.Scenario, func(t *testing.T) {
			bolt.DBPath = ".."
			DATABASE = testcase.Database
			err := SavePackagesToDisk("Test", "../encoder/test_"+testcase.Database+"/")
			if err != nil {
				t.Error(err)
			}
		})
	}
}

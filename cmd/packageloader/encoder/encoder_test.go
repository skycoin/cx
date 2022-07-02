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
			Scenario: "Test with Redis database",
			Database: "redis",
		},
		{
			Scenario: "Test with Bolt database",
			Database: "bolt",
		},
	}
	for _, testcase := range tests {
		t.Run(testcase.Scenario, func(t *testing.T) {
			bolt.DBPath = ".."
			err := SavePackagesToDisk("Test", "../encoder/test_"+testcase.Database+"/", testcase.Database)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

package graph

import (
	"fmt"
	"testing"

	"github.com/skycoin/cx/cmd/packageloader/bolt"
)

func TestTree(t *testing.T) {
	tests := []struct {
		Scenario       string
		Database       string
		ExpectedResult string
	}{
		{
			Scenario:       "Test with Redis database",
			Database:       "redis",
			ExpectedResult: "id=0, module=main, imports=4,1,2,\nid=1, module=testimport2, imports=2,3,\nid=2, module=testimport1, imports=5,\nid=3, module=testimport3, imports=2,\nid=4, module=os, imports=,\nid=5, module=gl, imports=,\n",
		},
		{
			Scenario:       "Test with Bolt database",
			Database:       "bolt",
			ExpectedResult: "id=0, module=main, imports=4,1,2,\nid=1, module=testimport2, imports=2,3,\nid=2, module=testimport1, imports=5,\nid=3, module=testimport3, imports=2,\nid=4, module=os, imports=,\nid=5, module=gl, imports=,\n",
		},
	}
	for _, testcase := range tests {
		t.Run(testcase.Scenario, func(t *testing.T) {
			bolt.DBPath = ".."
			output, err := GetImportGraph("TestTree", testcase.Database)
			if err != nil {
				t.Error(err)
			}

			if fmt.Sprintf("%q", output) != fmt.Sprintf("%q", testcase.ExpectedResult) {
				t.Error("The string produced by Tree does not match the expected result:")
				t.Errorf("%q", output)
				t.Errorf("%q", testcase.ExpectedResult)
			}
		})
	}
}

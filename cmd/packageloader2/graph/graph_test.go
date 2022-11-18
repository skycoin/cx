package graph

import (
	"fmt"
	"testing"

	"github.com/skycoin/cx/cmd/packageloader/bolt"
)

// Test currently fails because the graph is dependent on the order packages are added
// since loader adds packages concurrently the ids and order are different each time
// but the package and imports are correct
func TestTree(t *testing.T) {
	tests := []struct {
		Scenario       string
		Database       string
		ExpectedResult string
	}{
		{
			Scenario:       "Test with Redis database",
			Database:       "redis",
			ExpectedResult: "id=0, module=main, imports=4,2,1,\nid=1, module=testimport1, imports=5,\nid=2, module=testimport2, imports=1,3,\nid=3, module=testimport3, imports=1,\nid=4, module=os, imports=,\nid=5, module=gl, imports=,\n",
		},
		{
			Scenario:       "Test with Bolt database",
			Database:       "bolt",
			ExpectedResult: "id=0, module=main, imports=4,2,1,\nid=1, module=testimport1, imports=5,\nid=2, module=testimport2, imports=1,3,\nid=3, module=testimport3, imports=1,\nid=4, module=os, imports=,\nid=5, module=gl, imports=,\n",
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

package tree

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
			ExpectedResult: "main\n|--os\n|--testimport2\n|  |--testimport1\n|  |  `--gl\n|  `--testimport3\n|     `--testimport1\n|        `--gl\n`--testimport1\n   `--gl\n",
		},
		{
			Scenario:       "Test with Bolt database",
			Database:       "bolt",
			ExpectedResult: "main\n|--os\n|--testimport2\n|  |--testimport1\n|  |  `--gl\n|  `--testimport3\n|     `--testimport1\n|        `--gl\n`--testimport1\n   `--gl\n",
		},
	}
	for _, testcase := range tests {
		t.Run(testcase.Scenario, func(t *testing.T) {
			bolt.DBPath = ".."
			output, err := GetImportTree("TestTree", testcase.Database)
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

package encoder

import (
	"testing"

	"github.com/skycoin/cx/cmd/packageloader/bolt"
)

func TestSavePackage(t *testing.T) {
	for _, v := range []string{"redis", "bolt"} {
		bolt.DBPath = ".."
		DATABASE = v
		err := SavePackagesToDisk("Test", "../encoder/test_"+v+"/")
		if err != nil {
			t.Error(err)
		}
	}
}

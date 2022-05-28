package encoder

import (
	"testing"

	"github.com/skycoin/cx/cmd/packageloader/bolt"
)

func TestSavePackageRedis(t *testing.T) {
	bolt.DBPath = ".."
	DATABASE = "redis"
	err := SavePackagesToDisk("Test", "../encoder/test_redis/")
	if err != nil {
		t.Error(err)
	}
}

func TestSavePackageBolt(t *testing.T) {
	bolt.DBPath = ".."
	DATABASE = "bolt"
	err := SavePackagesToDisk("Test", "../encoder/test_bolt/")
	if err != nil {
		t.Error(err)
	}
}

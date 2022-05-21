package encoder

import "testing"

func TestSavePackage(t *testing.T) {
	SavePackagesToDisk("Test", "../encoder/test/")
}

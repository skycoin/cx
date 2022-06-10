package loader

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"

	"github.com/skycoin/cx/cmd/packageloader/bolt"
	"github.com/skycoin/cx/cmd/packageloader/redis"
	"golang.org/x/crypto/blake2b"
)

type PackageList struct {
	Packages []string
}

func (pl PackageList) MarshalBinary() ([]byte, error) {
	return json.Marshal(pl)
}

func (pl *PackageList) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, pl)
}

// Encode a package and put it in the specified package list
func (packageList *PackageList) hashPackage(newPackage *Package, database string) error {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(newPackage)
	if err != nil {
		return err
	}
	h := blake2b.Sum512(buffer.Bytes())
	packageList.Packages = append(packageList.Packages, fmt.Sprintf("%x", h[:]))
	switch database {
	case "redis":
		redis.Add(fmt.Sprintf("%x", h[:]), *newPackage)
	case "bolt":
		value, err := newPackage.MarshalBinary()
		if err != nil {
			return err
		}
		bolt.Add(fmt.Sprintf("%x", h[:]), value)
	}
	return nil
}

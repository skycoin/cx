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

type Package struct {
	PackageName string
	Files       []string
}

func (p Package) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Package) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}

// Encode a file and put it in the specified package
func (newPackage *Package) appendFile(newFile *File, database string) error {
	hash, err := newFile.getHash()
	if err != nil {
		return err
	}
	newPackage.Files = append(newPackage.Files, fmt.Sprintf("%x", hash[:]))
	err = newFile.saveToDatabase(hash, database)
	if err != nil {
		return err
	}
	return nil
}

func (newPackage *Package) getHash() ([64]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(newPackage)
	if err != nil {
		return [64]byte{}, err
	}
	return blake2b.Sum512(buffer.Bytes()), err
}

func (newPackage *Package) saveToDatabase(hash [64]byte, database string) error {
	switch database {
	case "redis":
		redis.Add(fmt.Sprintf("%x", hash[:]), *newPackage)
	case "bolt":
		value, err := newPackage.MarshalBinary()
		if err != nil {
			return err
		}
		bolt.Add(fmt.Sprintf("%x", hash[:]), value)
	}
	PackageHashMap[newPackage.PackageName] = fmt.Sprintf("%x", hash[:])
	return nil
}

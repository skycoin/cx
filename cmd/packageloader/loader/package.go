package loader

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"

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
func (newPackage *Package) hashFile(newFile *File) error {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(newFile)
	if err != nil {
		return err
	}
	h := blake2b.Sum512(buffer.Bytes())

	newPackage.Files = append(newPackage.Files, fmt.Sprintf("%x", h[:]))
	switch DATABASE {
	case "redis":
		redis.Add(fmt.Sprintf("%x", h[:]), *newFile)
	case "bolt":
		value, err := newFile.MarshalBinary()
		log.Print(value)
		if err != nil {
			log.Fatal(err)
		}
		bolt.Add(fmt.Sprintf("%x", h[:]), value)
	}
	return nil
}

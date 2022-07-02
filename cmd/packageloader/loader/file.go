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

type File struct {
	FileName   string
	Length     uint32
	Content    []byte
	Blake2Hash string
}

func (f File) MarshalBinary() ([]byte, error) {
	return json.Marshal(f)
}

func (f *File) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, f)
}

func (newFile *File) getHash() ([64]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(newFile)
	if err != nil {
		return [64]byte{}, err
	}
	return blake2b.Sum512(buffer.Bytes()), err
}

func (newFile *File) saveToDatabase(hash [64]byte, database string) error {
	switch database {
	case "redis":
		redis.Add(fmt.Sprintf("%x", hash[:]), *newFile)
	case "bolt":
		value, err := newFile.MarshalBinary()
		if err != nil {
			return err
		}
		bolt.Add(fmt.Sprintf("%x", hash[:]), value)
	}
	FileHashMap[newFile.FileName] = fmt.Sprintf("%x", hash[:])
	return nil
}

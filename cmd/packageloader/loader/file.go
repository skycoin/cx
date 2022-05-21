package loader

import "encoding/json"

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

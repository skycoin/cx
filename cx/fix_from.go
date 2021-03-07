fix_from.go

package cxcore

import(
	//"github.com/skycoin/skycoin/src/cipher/encoder"
)

//TODO:
//direct substitute these functions, instead of calling switch statement

func leUint16(b []byte) uint16 { return uint16(b[0]) | uint16(b[1])<<8 }

func lePutUint16(b []byte, v uint16) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
}

func leUint32(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

func lePutUint32(b []byte, v uint32) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
}

func leUint64(b []byte) uint64 {
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}

func lePutUint64(b []byte, v uint64) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 32)
	b[5] = byte(v >> 40)
	b[6] = byte(v >> 48)
	b[7] = byte(v >> 56)
}

func SerializeAtomic(data interface{}) []byte {
	var b [8]byte

	switch v := data.(type) {
	case bool:
		if v {
			b[0] = 1
		} else {
			b[0] = 0
		}
		return b[:1]
	case int8:
		b[0] = byte(v)
		return b[:1]
	case uint8:
		b[0] = v
		return b[:1]
	case int16:
		lePutUint16(b[:2], uint16(v))
		return b[:2]
	case uint16:
		lePutUint16(b[:2], v)
		return b[:2]
	case int32:
		lePutUint32(b[:4], uint32(v))
		return b[:4]
	case uint32:
		lePutUint32(b[:4], v)
		return b[:4]
	case int64:
		lePutUint64(b[:8], uint64(v))
		return b[:8]
	case uint64:
		lePutUint64(b[:8], v)
		return b[:8]
	default:
		log.Panic("SerializeAtomic unhandled type")
		return nil
	}
}


// FromI8 ...
func FromI8(in int8) []byte {
	//Serialize Atomic uses switch! Use serialize int8 directly
	//copy code over from encoder
	return SerializeAtomic(in)
}

// FromI16 ...
func FromI16(in int16) []byte {
	//Serialize Atomic uses switch! Use serialize int16 directly
	//copy code over
	return SerializeAtomic(in)
}

// FromI32 ...
func FromI32(in int32) []byte {
	//Serialize Atomic uses switch! Use serialize int32 directly
	//copy code over
	return SerializeAtomic(in)
}

// FromI64 ...
func FromI64(in int64) []byte {
	//Serialize Atomic uses switch! Use serialize int64 directly
	//copy code over
	return SerializeAtomic(in)
}

// FromUI8 ...
func FromUI8(in uint8) []byte {
	return SerializeAtomic(in)
}

// FromUI16 ...
func FromUI16(in uint16) []byte {
	return SerializeAtomic(in)
}

// FromUI32 ...
func FromUI32(in uint32) []byte {
	return SerializeAtomic(in)
}

// FromUI64 ...
func FromUI64(in uint64) []byte {
	return SerializeAtomic(in)
}

// FromF32 ...
func FromF32(in float32) []byte {
	return FromUI32(math.Float32bits(in))
}

// FromF64 ...
func FromF64(in float64) []byte {
	return FromUI64(math.Float64bits(in))
}



// FromStr ...
func FromStr(in string) []byte {
	return encoder.Serialize(in)
}

// FromBool ...
func FromBool(in bool) []byte {
	if in {
		return []byte{1}
	}
	return []byte{0}

}
package ameda

import "reflect"

// InitPointer initializes nil pointer with zero value.
func InitPointer(v reflect.Value) (done bool) {
	for {
		kind := v.Kind()
		if kind == reflect.Interface {
			v = v.Elem()
			continue
		}
		if kind != reflect.Ptr {
			return true
		}
		u := v.Elem()
		if u.IsValid() {
			v = u
			continue
		}
		if !v.CanSet() {
			return false
		}
		v2 := reflect.New(v.Type().Elem())
		v.Set(v2)
		v = v.Elem()
	}
}

// InitString initializes empty string pointer with def.
func InitString(p *string, def string) (done bool) {
	if p == nil {
		return false
	}
	if *p == "" {
		*p = def
	}
	return true
}

// InitBool initializes false bool pointer with def.
func InitBool(p *bool, def bool) (done bool) {
	if p == nil {
		return false
	}
	if *p == false {
		*p = def
	}
	return true
}

// InitByte initializes zero byte pointer with def.
func InitByte(p *byte, def byte) (done bool) {
	if p == nil {
		return false
	}
	if *p == 0 {
		*p = def
	}
	return true
}

// InitInt initializes zero int pointer with def.
func InitInt(p *int, def int) (done bool) {
	if p == nil {
		return false
	}
	if *p == 0 {
		*p = def
	}
	return true
}

// InitInt8 initializes zero int8 pointer with def.
func InitInt8(p *int8, def int8) (done bool) {
	if p == nil {
		return false
	}
	if *p == 0 {
		*p = def
	}
	return true
}

// InitInt16 initializes zero int16 pointer with def.
func InitInt16(p *int16, def int16) (done bool) {
	if p == nil {
		return false
	}
	if *p == 0 {
		*p = def
	}
	return true
}

// InitInt32 initializes zero int32 pointer with def.
func InitInt32(p *int32, def int32) (done bool) {
	if p == nil {
		return false
	}
	if *p == 0 {
		*p = def
	}
	return true
}

// InitInt64 initializes zero int64 pointer with def.
func InitInt64(p *int64, def int64) (done bool) {
	if p == nil {
		return false
	}
	if *p == 0 {
		*p = def
	}
	return true
}

// InitUint initializes zero uint pointer with def.
func InitUint(p *uint, def uint) (done bool) {
	if p == nil {
		return false
	}
	if *p == 0 {
		*p = def
	}
	return true
}

// InitUint8 initializes zero uint8 pointer with def.
func InitUint8(p *uint8, def uint8) (done bool) {
	if p == nil {
		return false
	}
	if *p == 0 {
		*p = def
	}
	return true
}

// InitUint16 initializes zero uint16 pointer with def.
func InitUint16(p *uint16, def uint16) (done bool) {
	if p == nil {
		return false
	}
	if *p == 0 {
		*p = def
	}
	return true
}

// InitUint32 initializes zero uint32 pointer with def.
func InitUint32(p *uint32, def uint32) (done bool) {
	if p == nil {
		return false
	}
	if *p == 0 {
		*p = def
	}
	return true
}

// InitUint64 initializes zero uint64 pointer with def.
func InitUint64(p *uint64, def uint64) (done bool) {
	if p == nil {
		return false
	}
	if *p == 0 {
		*p = def
	}
	return true
}

// InitFloat32 initializes zero float32 pointer with def.
func InitFloat32(p *float32, def float32) (done bool) {
	if p == nil {
		return false
	}
	if *p == 0 {
		*p = def
	}
	return true
}

// InitFloat64 initializes zero float64 pointer with def.
func InitFloat64(p *float64, def float64) (done bool) {
	if p == nil {
		return false
	}
	if *p == 0 {
		*p = def
	}
	return true
}

package interpreted

import (
	"fmt"
	"errors"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func castToStr (arg *CXArgument, expr *CXExpression, call *CXCall) error {
	strTyp := "str"
	switch arg.Typ {
	case "[]byte":
		assignOutput(0, *arg.Value, strTyp, expr, call)
	case "bool":
		var val int32
		encoder.DeserializeAtomic(*arg.Value, &val)
		var output []byte
		if val == 1 {
			output = encoder.Serialize("true")
		} else {
			output = encoder.Serialize("false")
		}
		assignOutput(0, output, strTyp, expr, call)
	case "byte":
		var val byte
		encoder.DeserializeRaw(*arg.Value, &val)
		output := encoder.Serialize(fmt.Sprintf("%d", val))
		assignOutput(0, output, strTyp, expr, call)
	case "i32":
		var val int32
		encoder.DeserializeAtomic(*arg.Value, &val)
		output := encoder.Serialize(fmt.Sprintf("%d", val))
		assignOutput(0, output, strTyp, expr, call)
	case "i64":
		var val int64
		encoder.DeserializeRaw(*arg.Value, &val)
		output := encoder.Serialize(fmt.Sprintf("%d", val))
		assignOutput(0, output, strTyp, expr, call)
	case "f32":
		var val float32
		encoder.DeserializeRaw(*arg.Value, &val)
		output := encoder.Serialize(fmt.Sprintf("%f", val))
		assignOutput(0, output, strTyp, expr, call)
	case "f64":
		var val float64
		encoder.DeserializeRaw(*arg.Value, &val)
		output := encoder.Serialize(fmt.Sprintf("%f", val))
		assignOutput(0, output, strTyp, expr, call)
	default:
		return errors.New(fmt.Sprintf("castToStr: type '%s' can't be casted to type 'str'", arg.Typ))
	}
	return nil
}

func castToByteA (arg *CXArgument, expr *CXExpression, call *CXCall) error {
	byteATyp := "[]byte"
	switch arg.Typ {
	case "str":
		assignOutput(0, *arg.Value, byteATyp, expr, call)
		return nil
	case "[]byte":
		assignOutput(0, *arg.Value, byteATyp, expr, call)
		return nil
	case "[]i32":
		var val []int32
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]byte, len(val) + 4)
		copy(output, (*arg.Value)[:4])
		for i, n := range val {
			output[i+4] = byte(n)
		}

		assignOutput(0, output, byteATyp, expr, call)
		return nil
	case "[]i64":
		var val []int64
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]byte, len(val) + 4)
		copy(output, (*arg.Value)[:4])
		for i, n := range val {
			output[i+4] = byte(n)
		}

		assignOutput(0, output, byteATyp, expr, call)
		return nil
	case "[]f32":
		var val []float32
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]byte, len(val) + 4)
		copy(output, (*arg.Value)[:4])
		for i, n := range val {
			output[i+4] = byte(n)
		}
		
		assignOutput(0, output, byteATyp, expr, call)
		return nil
	case "[]f64":
		var val []float64
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]byte, len(val) + 4)
		copy(output, (*arg.Value)[:4])
		for i, n := range val {
			output[i+4] = byte(n)
		}
		
		assignOutput(0, output, byteATyp, expr, call)
		return nil
	default:
		return errors.New(fmt.Sprintf("castToByteA: type '%s' can't be casted to type '[]byte'", arg.Typ))
	}
}

func castToByte (arg *CXArgument, expr *CXExpression, call *CXCall) error {
	byteTyp := "byte"
	switch arg.Typ {
	case "byte":
		assignOutput(0, *arg.Value, byteTyp, expr, call)
		return nil
	case "i32":
		var val int32
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := []byte{byte(val)}

		assignOutput(0, newVal, byteTyp, expr, call)
		return nil
	case "i64":
		var val int64
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := []byte{byte(val)}

		assignOutput(0, newVal, byteTyp, expr, call)
		return nil
	case "f32":
		var val float32
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := []byte{byte(val)}

		assignOutput(0, newVal, byteTyp, expr, call)
		return nil
	case "f64":
		var val float64
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := []byte{byte(val)}

		assignOutput(0, newVal, byteTyp, expr, call)
		return nil
	default:
		return errors.New(fmt.Sprintf("castToByte: type '%s' can't be casted to type 'byte'", arg.Typ))
	}
}

func castToI32 (arg *CXArgument, expr *CXExpression, call *CXCall) error {
	i32Typ := "i32"
	switch arg.Typ {
	case "byte":
		val := (*arg.Value)[0]
		newVal := encoder.Serialize(int32(val))

		assignOutput(0, newVal, i32Typ, expr, call)
		return nil
	case "i32":
		assignOutput(0, *arg.Value, i32Typ, expr, call)
		return nil
	case "i64":
		var val int64
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(int32(val))

		assignOutput(0, newVal, i32Typ, expr, call)
		return nil
	case "f32":
		var val float32
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(int32(val))

		assignOutput(0, newVal, i32Typ, expr, call)
		return nil
	case "f64":
		var val float64
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(int32(val))

		assignOutput(0, newVal, i32Typ, expr, call)
		return nil
	default:
		return errors.New(fmt.Sprintf("castToI32: type '%s' can't be casted to type 'i32'", arg.Typ))
	}
}

func castToI64 (arg *CXArgument, expr *CXExpression, call *CXCall) error {
	i64Typ := "i64"
	switch arg.Typ {
	case "byte":
		val := (*arg.Value)[0]
		newVal := encoder.Serialize(int64(val))

		assignOutput(0, newVal, i64Typ, expr, call)
		return nil
	case "i32":
		var val int32
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(int64(val))

		assignOutput(0, newVal, i64Typ, expr, call)
		return nil
	case "i64":
		assignOutput(0, *arg.Value, i64Typ, expr, call)
		return nil
	case "f32":
		var val float32
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(int64(val))

		assignOutput(0, newVal, i64Typ, expr, call)
		return nil
	case "f64":
		var val float64
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(int64(val))

		assignOutput(0, newVal, i64Typ, expr, call)
		return nil
	default:
		return errors.New(fmt.Sprintf("castToI64: type '%s' can't be casted to type 'i64'", arg.Typ))
	}
}

func castToF32 (arg *CXArgument, expr *CXExpression, call *CXCall) error {
	f32Typ := "f32"
	switch arg.Typ {
	case "byte":
		val := (*arg.Value)[0]
		newVal := encoder.Serialize(float32(val))

		assignOutput(0, newVal, f32Typ, expr, call)
		return nil
	case "i32":
		var val int32
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(float32(val))

		assignOutput(0, newVal, f32Typ, expr, call)
		return nil
	case "i64":
		var val int64
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(float32(val))

		assignOutput(0, newVal, f32Typ, expr, call)
		return nil
	case "f32":
		assignOutput(0, *arg.Value, f32Typ, expr, call)
		return nil
	case "f64":
		var val float64
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(float32(val))

		assignOutput(0, newVal, f32Typ, expr, call)
		return nil
	default:
		return errors.New(fmt.Sprintf("castToF32: type '%s' can't be casted to type 'f32'", arg.Typ))
	}
}

func castToF64 (arg *CXArgument, expr *CXExpression, call *CXCall) error {
	f64Typ := "f64"
	switch arg.Typ {
	case "byte":
		val := (*arg.Value)[0]
		newVal := encoder.Serialize(float64(val))

		assignOutput(0, newVal, f64Typ, expr, call)
		return nil
	case "i32":
		var val int32
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(float64(val))

		assignOutput(0, newVal, f64Typ, expr, call)
		return nil
	case "i64":
		var val int64
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(float64(val))

		assignOutput(0, newVal, f64Typ, expr, call)
		return nil
	case "f32":
		var val float32
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(float64(val))

		assignOutput(0, newVal, f64Typ, expr, call)
		return nil
	case "f64":
		assignOutput(0, *arg.Value, f64Typ, expr, call)
		return nil
	default:
		return errors.New(fmt.Sprintf("castToF64: type '%s' can't be casted to type 'f64'", arg.Typ))
	}
}

func castToI32A (arg *CXArgument, expr *CXExpression, call *CXCall) error {
	i32ATyp := "[]i32"
	switch arg.Typ {
	case "[]byte":
		var val []byte
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]int32, len(val))
		for i, n := range val {
			output[i] = int32(n)
		}
		newVal := encoder.Serialize(output)

		assignOutput(0, newVal, i32ATyp, expr, call)
		return nil
	case "[]i32":
		assignOutput(0, *arg.Value, i32ATyp, expr, call)
		return nil
	case "[]i64":
		var val []int64
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]int32, len(val))
		for i, n := range val {
			output[i] = int32(n)
		}		
		newVal := encoder.Serialize(output)

		assignOutput(0, newVal, i32ATyp, expr, call)
		return nil
	case "[]f32":
		var val []float32
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]int32, len(val))
		for i, n := range val {
			output[i] = int32(n)
		}
		newVal := encoder.Serialize(output)

		assignOutput(0, newVal, i32ATyp, expr, call)
		return nil
	case "[]f64":
		var val []float64
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]int32, len(val))
		for i, n := range val {
			output[i] = int32(n)
		}
		newVal := encoder.Serialize(output)

		assignOutput(0, newVal, i32ATyp, expr, call)
		return nil
	default:
		return errors.New(fmt.Sprintf("castToI32A: type '%s' can't be casted to type '[]i32'", arg.Typ))
	}
}

func castToI64A (arg *CXArgument, expr *CXExpression, call *CXCall) error {
	i64ATyp := "[]i64"
	switch arg.Typ {
	case "[]byte":
		var val []byte
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]int64, len(val))
		for i, n := range val {
			output[i] = int64(n)
		}
		newVal := encoder.Serialize(output)

		assignOutput(0, newVal, i64ATyp, expr, call)
		return nil
	case "[]i32":
		var val []int32
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]int64, len(val))
		for i, n := range val {
			output[i] = int64(n)
		}
		newVal := encoder.Serialize(output)

		assignOutput(0, newVal, i64ATyp, expr, call)
		return nil
	case "[]i64":
		assignOutput(0, *arg.Value, i64ATyp, expr, call)
		return nil
	case "[]f32":
		var val []float32
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]int64, len(val))
		for i, n := range val {
			output[i] = int64(n)
		}
		newVal := encoder.Serialize(output)

		assignOutput(0, newVal, i64ATyp, expr, call)
		return nil
	case "[]f64":
		var val []float64
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]int64, len(val))
		for i, n := range val {
			output[i] = int64(n)
		}
		newVal := encoder.Serialize(output)

		assignOutput(0, newVal, i64ATyp, expr, call)
		return nil
	default:
		return errors.New(fmt.Sprintf("castToI64A: type '%s' can't be casted to type '[]i64'", arg.Typ))
	}
}

func castToF32A (arg *CXArgument, expr *CXExpression, call *CXCall) error {
	f32ATyp := "[]f32"
	switch arg.Typ {
	case "[]byte":
		var val []byte
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]float32, len(val))
		for i, n := range val {
			output[i] = float32(n)
		}
		newVal := encoder.Serialize(output)

		assignOutput(0, newVal, f32ATyp, expr, call)
		return nil
	case "[]i32":
		var val []int32
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]float32, len(val))
		for i, n := range val {
			output[i] = float32(n)
		}
		newVal := encoder.Serialize(output)

		assignOutput(0, newVal, f32ATyp, expr, call)
		return nil
	case "[]i64":
		var val []int64
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]float32, len(val))
		for i, n := range val {
			output[i] = float32(n)
		}
		newVal := encoder.Serialize(output)

		assignOutput(0, newVal, f32ATyp, expr, call)
		return nil
	case "[]f32":
		assignOutput(0, *arg.Value, f32ATyp, expr, call)
		return nil
	case "[]f64":
		var val []float64
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]float32, len(val))
		for i, n := range val {
			output[i] = float32(n)
		}
		
		newVal := encoder.Serialize(output)

		assignOutput(0, newVal, f32ATyp, expr, call)
		return nil
	default:
		return errors.New(fmt.Sprintf("castToF32A: type '%s' can't be casted to type '[]f32'", arg.Typ))
	}
}

func castToF64A (arg *CXArgument, expr *CXExpression, call *CXCall) error {
	f64ATyp := "[]f64"
	switch arg.Typ {
	case "[]byte":
		var val []byte
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]float64, len(val))
		for i, n := range val {
			output[i] = float64(n)
		}
		newVal := encoder.Serialize(output)

		assignOutput(0, newVal, f64ATyp, expr, call)
		return nil
	case "[]i32":
		var val []int32
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]float64, len(val))
		for i, n := range val {
			output[i] = float64(n)
		}
		newVal := encoder.Serialize(output)

		assignOutput(0, newVal, f64ATyp, expr, call)
		return nil
	case "[]i64":
		var val []int64
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]float64, len(val))
		for i, n := range val {
			output[i] = float64(n)
		}
		newVal := encoder.Serialize(output)

		assignOutput(0, newVal, f64ATyp, expr, call)
		return nil
	case "[]f32":
		var val []float32
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]float64, len(val))
		for i, n := range val {
			output[i] = float64(n)
		}
		newVal := encoder.Serialize(output)

		assignOutput(0, newVal, f64ATyp, expr, call)
		return nil
	case "[]f64":
		assignOutput(0, *arg.Value, f64ATyp, expr, call)
		return nil
	default:
		return errors.New(fmt.Sprintf("castToF64A: type '%s' can't be casted to type '[]f64'", arg.Typ))
	}
}

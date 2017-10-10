package base

import (
	"fmt"
	"errors"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func castToStr (arg *CXArgument, expr *CXExpression, call *CXCall) error {
	switch arg.Typ {
	case "[]byte":
		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = arg.Value
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, arg.Value, "str"))
		return nil
	default:
		return errors.New(fmt.Sprintf("castToStr: type '%s' can't be casted to type 'str'", arg.Typ))
	}
}

func castToByteA (arg *CXArgument, expr *CXExpression, call *CXCall) error {
	byteATyp := "[]byte"
	switch arg.Typ {
	case "str":
		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = arg.Value
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, arg.Value, byteATyp))
		return nil
	case "[]byte":
		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = arg.Value
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, arg.Value, byteATyp))
		return nil
	case "[]i32":
		var val []int32
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]byte, len(val))
		for i, n := range val {
			output[i] = byte(n)
		}

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &output
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &output, byteATyp))
		return nil
	case "[]i64":
		var val []int64
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]byte, len(val))
		for i, n := range val {
			output[i] = byte(n)
		}
		
		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &output
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &output, byteATyp))
		return nil
	case "[]f32":
		var val []float32
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]byte, len(val))
		for i, n := range val {
			output[i] = byte(n)
		}
		
		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &output
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &output, byteATyp))
		return nil
	case "[]f64":
		var val []float64
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]byte, len(val))
		for i, n := range val {
			output[i] = byte(n)
		}
		
		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &output
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &output, byteATyp))
		return nil
	default:
		return errors.New(fmt.Sprintf("castToByteA: type '%s' can't be casted to type '[]byte'", arg.Typ))
	}
}

func castToByte (arg *CXArgument, expr *CXExpression, call *CXCall) error {
	byteTyp := "byte"
	switch arg.Typ {
	case "byte":
		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = arg.Value
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, arg.Value, byteTyp))
		return nil
	case "i32":
		var val int32
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := []byte{byte(val)}

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &newVal
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &newVal, byteTyp))
		return nil
	case "i64":
		var val int64
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := []byte{byte(val)}

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &newVal
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &newVal, byteTyp))
		return nil
	case "f32":
		var val float32
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := []byte{byte(val)}

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &newVal
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &newVal, byteTyp))
		return nil
	case "f64":
		var val float64
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := []byte{byte(val)}

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &newVal
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &newVal, byteTyp))
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

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &newVal
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &newVal, i32Typ))
		return nil
	case "i32":
		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = arg.Value
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, arg.Value, i32Typ))
		return nil
	case "i64":
		var val int64
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(int32(val))

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &newVal
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &newVal, i32Typ))
		return nil
	case "f32":
		var val float32
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(int32(val))

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &newVal
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &newVal, i32Typ))
		return nil
	case "f64":
		var val float64
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(int32(val))

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &newVal
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &newVal, i32Typ))
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

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &newVal
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &newVal, i64Typ))
		return nil
	case "i32":
		var val int32
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(int64(val))

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &newVal
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &newVal, i64Typ))
		return nil
	case "i64":
		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = arg.Value
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, arg.Value, i64Typ))
		return nil
	case "f32":
		var val float32
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(int64(val))

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &newVal
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &newVal, i64Typ))
		return nil
	case "f64":
		var val float64
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(int64(val))

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &newVal
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &newVal, i64Typ))
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

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &newVal
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &newVal, f32Typ))
		return nil
	case "i32":
		var val int32
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(float32(val))

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &newVal
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &newVal, f32Typ))
		return nil
	case "i64":
		var val int64
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(float32(val))

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &newVal
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &newVal, f32Typ))
		return nil
	case "f32":
		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = arg.Value
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, arg.Value, f32Typ))
		return nil
	case "f64":
		var val float64
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(float32(val))

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &newVal
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &newVal, f32Typ))
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

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &newVal
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &newVal, f64Typ))
		return nil
	case "i32":
		var val int32
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(float64(val))

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &newVal
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &newVal, f64Typ))
		return nil
	case "i64":
		var val int64
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(float64(val))

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &newVal
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &newVal, f64Typ))
		return nil
	case "f32":
		var val float32
		encoder.DeserializeRaw(*arg.Value, &val)
		newVal := encoder.Serialize(float64(val))

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &newVal
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &newVal, f64Typ))
		return nil
	case "f64":
		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = arg.Value
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, arg.Value, f64Typ))
		return nil
	default:
		return errors.New(fmt.Sprintf("castToF64: type '%s' can't be casted to type 'f64'", arg.Typ))
	}
}

func castToI32A (arg *CXArgument, expr *CXExpression, call *CXCall) error {
	i32ATyp := "[]i32"
	switch arg.Typ {
	case "[]byte":
		val := *arg.Value

		output := make([]int32, len(val))
		for i, n := range val {
			output[i] = int32(n)
		}
		newVal := encoder.Serialize(output)
		
		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &newVal
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &newVal, i32ATyp))
		return nil
	case "[]i32":
		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = arg.Value
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, arg.Value, i32ATyp))
		return nil
	case "[]i64":
		var val []int64
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]int32, len(val))
		for i, n := range val {
			output[i] = int32(n)
		}		
		newVal := encoder.Serialize(output)

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &newVal
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &newVal, i32ATyp))
		return nil
	case "[]f32":
		var val []float32
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]int32, len(val))
		for i, n := range val {
			output[i] = int32(n)
		}
		newVal := encoder.Serialize(output)

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &newVal
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &newVal, i32ATyp))
		return nil
	case "[]f64":
		var val []float64
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]int32, len(val))
		for i, n := range val {
			output[i] = int32(n)
		}
		newVal := encoder.Serialize(output)

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &newVal
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &newVal, i32ATyp))
		return nil
	default:
		return errors.New(fmt.Sprintf("castToI32A: type '%s' can't be casted to type '[]i32'", arg.Typ))
	}
}

func castToI64A (arg *CXArgument, expr *CXExpression, call *CXCall) error {
	i64ATyp := "[]i64"
	switch arg.Typ {
	case "[]byte":
		val := *arg.Value

		output := make([]int64, len(val))
		for i, n := range val {
			output[i] = int64(n)
		}
		newVal := encoder.Serialize(output)

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &newVal
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &newVal, i64ATyp))
		return nil
	case "[]i32":
		var val []int32
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]int64, len(val))
		for i, n := range val {
			output[i] = int64(n)
		}
		newVal := encoder.Serialize(output)

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &newVal
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &newVal, i64ATyp))
		return nil
	case "[]i64":
		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = arg.Value
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, arg.Value, i64ATyp))
		return nil
	case "[]f32":
		var val []float32
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]int64, len(val))
		for i, n := range val {
			output[i] = int64(n)
		}
		newVal := encoder.Serialize(output)

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &newVal
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &newVal, i64ATyp))
		return nil
	case "[]f64":
		var val []float64
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]int64, len(val))
		for i, n := range val {
			output[i] = int64(n)
		}
		newVal := encoder.Serialize(output)

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &newVal
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &newVal, i64ATyp))
		return nil
	default:
		return errors.New(fmt.Sprintf("castToI64A: type '%s' can't be casted to type '[]i64'", arg.Typ))
	}
}

func castToF32A (arg *CXArgument, expr *CXExpression, call *CXCall) error {
	f32ATyp := "[]f32"
	switch arg.Typ {
	case "[]byte":
		val := *arg.Value

		output := make([]float32, len(val))
		for i, n := range val {
			output[i] = float32(n)
		}
		newVal := encoder.Serialize(output)

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &newVal
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &newVal, f32ATyp))
		return nil
	case "[]i32":
		var val []int32
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]float32, len(val))
		for i, n := range val {
			output[i] = float32(n)
		}
		newVal := encoder.Serialize(output)

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &newVal
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &newVal, f32ATyp))
		return nil
	case "[]i64":
		var val []int64
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]float32, len(val))
		for i, n := range val {
			output[i] = float32(n)
		}
		newVal := encoder.Serialize(output)

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &newVal
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &newVal, f32ATyp))
		return nil
	case "[]f32":
		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = arg.Value
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, arg.Value, f32ATyp))
		return nil
	case "[]f64":
		var val []float64
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]float32, len(val))
		for i, n := range val {
			output[i] = float32(n)
		}
		
		newVal := encoder.Serialize(output)

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &newVal
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &newVal, f32ATyp))
		return nil
	default:
		return errors.New(fmt.Sprintf("castToF32A: type '%s' can't be casted to type '[]f32'", arg.Typ))
	}
}

func castToF64A (arg *CXArgument, expr *CXExpression, call *CXCall) error {
	f64ATyp := "[]f64"
	switch arg.Typ {
	case "[]byte":
		val := *arg.Value

		output := make([]float64, len(val))
		for i, n := range val {
			output[i] = float64(n)
		}
		newVal := encoder.Serialize(output)

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &newVal
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &newVal, f64ATyp))
		return nil
	case "[]i32":
		var val []int32
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]float64, len(val))
		for i, n := range val {
			output[i] = float64(n)
		}
		newVal := encoder.Serialize(output)

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &newVal
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &newVal, f64ATyp))
		return nil
	case "[]i64":
		var val []int64
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]float64, len(val))
		for i, n := range val {
			output[i] = float64(n)
		}
		newVal := encoder.Serialize(output)

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &newVal
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &newVal, f64ATyp))
		return nil
	case "[]f32":
		var val []float32
		encoder.DeserializeRaw(*arg.Value, &val)

		output := make([]float64, len(val))
		for i, n := range val {
			output[i] = float64(n)
		}
		newVal := encoder.Serialize(output)

		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = &newVal
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, &newVal, f64ATyp))
		return nil
	case "[]f64":
		for _, def := range call.State {
			if def.Name == expr.OutputNames[0].Name {
				def.Value = arg.Value
				return nil
			}
		}
		
		call.State = append(call.State, MakeDefinition(expr.OutputNames[0].Name, arg.Value, f64ATyp))
		return nil
	default:
		return errors.New(fmt.Sprintf("castToF64A: type '%s' can't be casted to type '[]f64'", arg.Typ))
	}
}

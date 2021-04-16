package ast

import "github.com/skycoin/skycoin/src/cipher/encoder"

// SerializeCXProgramV2 translates cx program to slice of bytes that we can save.
// These slice of bytes can then be deserialize in the future and
// be translated back to cx program.
func SerializeCXProgramV2(prgrm *CXProgram, includeMemory bool) (b []byte) {
	s := SerializedCXProgram{}
	initSerialization(prgrm, &s, includeMemory)

	// serialize cx program's packages,
	// structs, functions, etc.
	serializeCXProgramElements(prgrm, &s)

	// serialize cx program's program
	serializeProgram(prgrm, &s)

	// assign cx program's offsets
	assignSerializedCXProgramOffset(&s)

	convertSerializedCXProgramMapsToKVPairs(&s)

	// serializing everything
	b = encoder.Serialize(s)

	return b
}

func convertMapToKVPairs(inputMap map[string]int64) []KeyValuePair {
	kvPairs := make([]KeyValuePair, len(inputMap))
	index := 0
	for k, v := range inputMap {
		kvPairs[index].Key = k
		kvPairs[index].Value = v
		index++
	}

	return kvPairs
}

func convertKVPairsToMap(inputKVPairs []KeyValuePair) map[string]int64 {
	outputMap := make(map[string]int64)
	for _, kv := range inputKVPairs {
		outputMap[kv.Key] = kv.Value
	}

	return outputMap
}

func convertSerializedCXProgramMapsToKVPairs(s *SerializedCXProgram) {
	s.PackagesMapKV = convertMapToKVPairs(s.PackagesMap)
	s.PackagesMap = make(map[string]int64)

	s.StructsMapKV = convertMapToKVPairs(s.StructsMap)
	s.StructsMap = make(map[string]int64)

	s.FunctionsMapKV = convertMapToKVPairs(s.FunctionsMap)
	s.FunctionsMap = make(map[string]int64)

	s.StringsMapKV = convertMapToKVPairs(s.StringsMap)
	s.StringsMap = make(map[string]int64)
}

func convertSerializedCXProgramKVPairsToMaps(s *SerializedCXProgram) {
	s.PackagesMap = make(map[string]int64)
	s.PackagesMap = convertKVPairsToMap(s.PackagesMapKV)

	s.StructsMap = make(map[string]int64)
	s.StructsMap = convertKVPairsToMap(s.StructsMapKV)

	s.FunctionsMap = make(map[string]int64)
	s.FunctionsMap = convertKVPairsToMap(s.FunctionsMapKV)

	s.StringsMap = make(map[string]int64)
	s.StringsMap = convertKVPairsToMap(s.StringsMapKV)
}

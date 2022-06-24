package types

//"fmt"

func AllocWrite_str_data(prgrm interface{}, memory []byte, str string) Pointer {
	return AllocWrite_obj_data(prgrm, memory, []byte(str))
}

func Write_str_data(memory []byte, offset Pointer, value string) {
	bytes := []byte(value)
	size := Compute_obj_size(bytes)
	Write_obj_data(memory, offset, size, bytes)
}

func Read_str_data(memory []byte, offset Pointer) string {
	str := Read_obj_data(memory, offset)
	return string(str)
}

func Write_str(prgrm interface{}, memory []byte, offset Pointer, str string) {
	Write_obj(prgrm, memory, offset, []byte(str))
}

func Read_str(memory []byte, offset Pointer) string {
	str := string(Read_obj(memory, offset))
	return str
}

func Read_str_size(memory []byte, offset Pointer) Pointer {
	heapOffset := Read_ptr(memory, offset)
	if heapOffset > 0 && heapOffset.IsValid() {
		return Read_obj_size(memory, heapOffset) - OBJECT_HEADER_SIZE
	}
	return 0
}

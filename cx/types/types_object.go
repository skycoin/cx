package types

//"fmt"

const MARK_SIZE = Pointer(1)
const FORWARDING_ADDRESS_SIZE = POINTER_SIZE
const OBJECT_GC_HEADER_SIZE = MARK_SIZE + FORWARDING_ADDRESS_SIZE

const OBJECT_SIZE = POINTER_SIZE
const OBJECT_HEADER_SIZE = OBJECT_GC_HEADER_SIZE + OBJECT_SIZE

func Get_obj_header(memory []byte, offset Pointer) []byte {
	return memory[offset : offset+OBJECT_HEADER_SIZE]
}

func Get_obj_data(memory []byte, offset Pointer, size Pointer) []byte {
	offset += OBJECT_HEADER_SIZE
	return memory[offset : offset+size]
}

func Compute_obj_size(data []byte) Pointer {
	return Cast_int_to_ptr(len(data)) + OBJECT_HEADER_SIZE
}

func Read_obj_forwarding_address(memory []byte, offset Pointer) Pointer {
	return Read_ptr(memory, offset+MARK_SIZE)
}

func Read_obj_size(memory []byte, offset Pointer) Pointer {
	return Read_ptr(memory, offset+OBJECT_HEADER_SIZE-OBJECT_SIZE)
}

func Write_obj_size(memory []byte, offset Pointer, size Pointer) {
	Write_ptr(memory, offset+OBJECT_HEADER_SIZE-OBJECT_SIZE, size)
}

func Write_obj_mark(memory []byte, offset Pointer, mark byte) {
	Write_ui8(memory, offset, mark)
}

func Write_obj_forwarding_address(memory []byte, offset Pointer, address Pointer) {
	Write_ptr(memory, offset+MARK_SIZE, address)
}

func Make_obj(data []byte) []byte {
	size := Compute_obj_size(data)
	obj := make([]byte, size)
	//fmt.Printf("MAKE_OBJ SIZE %d, LEN_DATA %d, LEN_OBJ %d\n", size, len(data), len(obj))
	Write_obj_data(obj, 0, data)
	return obj
}

func AllocWrite_obj_data(prgrm interface{}, memory []byte, obj []byte) Pointer {
	heapOffset := Allocator(prgrm, Compute_obj_size(obj))
	Write_obj_data(memory, heapOffset, obj)
	return heapOffset
}

func Write_obj_data(memory []byte, offset Pointer, obj []byte) {
	size := Cast_int_to_ptr(len(obj))
	Write_obj_size(memory, offset, size)
	WriteSlice_byte(memory, offset+OBJECT_HEADER_SIZE, obj)
}

func Read_obj_data(memory []byte, offset Pointer) []byte {
	size := Read_obj_size(memory, offset)
	obj := GetSlice_byte(memory, offset+OBJECT_HEADER_SIZE, size)
	return obj
}

func Write_obj(prgrm interface{}, memory []byte, offset Pointer, obj []byte) {
	heapOffset := AllocWrite_obj_data(prgrm, memory, obj)
	Write_ptr(memory, offset, heapOffset)
}

func Read_obj(memory []byte, offset Pointer) []byte {
	heapOffset := Read_ptr(memory, offset)
	if heapOffset != 0 && heapOffset.IsValid() {
		obj := Read_obj_data(memory, heapOffset)
		return obj
	}
	return nil
}

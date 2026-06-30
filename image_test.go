package amx

import (
	"encoding/binary"

	"github.com/pawnkit/goamx/vm"
)

func debugTestAMX() []byte {
	data := publicTestAMX(vm.OP_HALT, 0)
	binary.LittleEndian.PutUint16(data[8:10], FlagDebug)

	chunk := make([]byte, 22)
	binary.LittleEndian.PutUint16(chunk[4:6], 0xf1ef)
	chunk[6], chunk[7] = 8, 11
	for offset := 10; offset <= 20; offset += 2 {
		binary.LittleEndian.PutUint16(chunk[offset:offset+2], 1)
	}

	appendUint32 := func(value uint32) {
		var raw [4]byte
		binary.LittleEndian.PutUint32(raw[:], value)
		chunk = append(chunk, raw[:]...)
	}
	appendInt16 := func(value uint16) {
		var raw [2]byte
		binary.LittleEndian.PutUint16(raw[:], value)
		chunk = append(chunk, raw[:]...)
	}
	appendName := func(value string) {
		chunk = append(chunk, value...)
		chunk = append(chunk, 0)
	}

	appendUint32(0)
	appendName("test.pwn")
	appendUint32(0)
	appendUint32(12)
	appendUint32(0)
	appendInt16(1)
	appendUint32(0)
	appendUint32(8)
	chunk = append(chunk, SymbolFunction, 0)
	appendInt16(1)
	appendName("test_public")
	appendInt16(2)
	appendUint32(3)
	appendInt16(1)
	appendName("Float")
	appendInt16(1)
	appendUint32(4)
	appendName("machine")
	appendInt16(2)
	appendInt16(1)
	appendName("state")

	binary.LittleEndian.PutUint32(chunk[0:4], uint32(len(chunk)))
	return append(data, chunk...)
}

func publicNativeAMX() []byte {
	const headerSize = 56
	publics := uint32(headerSize)
	natives := publics + 8
	libraries := natives + 8
	data := make([]byte, libraries)
	data = append(data, 31, 0)
	binary.LittleEndian.PutUint32(data[publics+4:publics+8], uint32(len(data)))
	data = append(data, "test_native"...)
	data = append(data, 0)
	binary.LittleEndian.PutUint32(data[natives+4:natives+8], uint32(len(data)))
	data = append(data, "increment"...)
	data = append(data, 0)
	for len(data)%4 != 0 {
		data = append(data, 0)
	}
	cod := uint32(len(data))
	for _, value := range []int32{int32(vm.OP_PUSH_C), 41, int32(vm.OP_SYSREQ_N), 0, 4, int32(vm.OP_HALT), 0} {
		var cell [4]byte
		binary.LittleEndian.PutUint32(cell[:], uint32(value))
		data = append(data, cell[:]...)
	}
	dat := uint32(len(data))
	binary.LittleEndian.PutUint32(data[0:4], dat)
	binary.LittleEndian.PutUint16(data[4:6], 0xf1e0)
	data[6], data[7] = 8, 11
	binary.LittleEndian.PutUint16(data[10:12], 8)
	binary.LittleEndian.PutUint32(data[12:16], cod)
	binary.LittleEndian.PutUint32(data[16:20], dat)
	binary.LittleEndian.PutUint32(data[20:24], dat)
	binary.LittleEndian.PutUint32(data[24:28], dat+256)
	binary.LittleEndian.PutUint32(data[32:36], publics)
	binary.LittleEndian.PutUint32(data[36:40], natives)
	binary.LittleEndian.PutUint32(data[40:44], libraries)
	binary.LittleEndian.PutUint32(data[44:48], libraries)
	binary.LittleEndian.PutUint32(data[48:52], libraries)
	binary.LittleEndian.PutUint32(data[52:56], libraries)
	return data
}

func publicTestAMX(cells ...any) []byte {
	const headerSize = 56
	publics := uint32(headerSize)
	natives := publics + 8
	data := make([]byte, natives)
	data = append(data, 31, 0)
	binary.LittleEndian.PutUint32(data[publics+4:publics+8], uint32(len(data)))
	data = append(data, "test_public"...)
	data = append(data, 0)
	for len(data)%4 != 0 {
		data = append(data, 0)
	}
	cod := uint32(len(data))
	for _, item := range cells {
		var value int32
		switch item := item.(type) {
		case vm.Opcode:
			value = int32(item)
		case int:
			value = int32(item)
		default:
			panic("unsupported test cell")
		}
		var cell [4]byte
		binary.LittleEndian.PutUint32(cell[:], uint32(value))
		data = append(data, cell[:]...)
	}
	dat := uint32(len(data))
	binary.LittleEndian.PutUint32(data[0:4], uint32(len(data)))
	binary.LittleEndian.PutUint16(data[4:6], 0xf1e0)
	data[6], data[7] = 8, 11
	binary.LittleEndian.PutUint16(data[10:12], 8)
	binary.LittleEndian.PutUint32(data[12:16], cod)
	binary.LittleEndian.PutUint32(data[16:20], dat)
	binary.LittleEndian.PutUint32(data[20:24], dat)
	binary.LittleEndian.PutUint32(data[24:28], dat+256)
	binary.LittleEndian.PutUint32(data[32:36], publics)
	binary.LittleEndian.PutUint32(data[36:40], natives)
	for offset := 40; offset <= 52; offset += 4 {
		binary.LittleEndian.PutUint32(data[offset:offset+4], natives)
	}
	return data
}

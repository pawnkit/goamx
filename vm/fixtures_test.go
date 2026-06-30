package vm

import (
	"encoding/binary"
	"testing"
)

func minimalAMX(t *testing.T) []byte {
	t.Helper()
	const (
		publics = uint32(headerSize)
		natives = publics + 2*8
		libs    = natives + 1*8
		names   = libs
	)
	data := make([]byte, names)
	putHeader(data, names)
	data = append(data, maxNameLength, 0)
	data = appendName(data, publics+4, "test_addition")
	data = appendName(data, publics+12, "helper")
	data = appendName(data, natives+4, "test_native")
	for len(data)%cellBytes != 0 {
		data = append(data, 0)
	}
	putHeader(data, uint32(len(data)))
	return data
}

func amxWithPubVarsAndTags(t *testing.T) []byte {
	t.Helper()
	const (
		publics = uint32(headerSize)
		natives = publics
		libs    = natives
		pubvars = libs
		tags    = pubvars + 8
		names   = tags + 8
	)
	data := make([]byte, names)
	putHeader(data, names)
	data = append(data, maxNameLength, 0)
	binary.LittleEndian.PutUint32(data[32:36], publics)
	binary.LittleEndian.PutUint32(data[36:40], natives)
	binary.LittleEndian.PutUint32(data[40:44], libs)
	binary.LittleEndian.PutUint32(data[44:48], pubvars)
	binary.LittleEndian.PutUint32(data[48:52], tags)
	binary.LittleEndian.PutUint32(data[52:56], names)
	data = appendName(data, pubvars+4, "global_value")
	data = appendName(data, tags+4, "Vehicle")
	for len(data)%cellBytes != 0 {
		data = append(data, 0)
	}
	putHeader(data, uint32(len(data)))
	binary.LittleEndian.PutUint32(data[32:36], publics)
	binary.LittleEndian.PutUint32(data[36:40], natives)
	binary.LittleEndian.PutUint32(data[40:44], libs)
	binary.LittleEndian.PutUint32(data[44:48], pubvars)
	binary.LittleEndian.PutUint32(data[48:52], tags)
	binary.LittleEndian.PutUint32(data[52:56], names)
	return data
}

func amxWithCode(t *testing.T, code []byte) []byte {
	return amxWithCodeAndNatives(t, code, nil)
}

func amxWithCodeAndNatives(t *testing.T, code []byte, nativeNames []string) []byte {
	t.Helper()
	publics := uint32(headerSize)
	natives := publics + 8
	libs := natives + uint32(len(nativeNames))*8
	data := make([]byte, libs)
	data = append(data, maxNameLength, 0)
	data = appendName(data, publics+4, "test_main")
	for i, name := range nativeNames {
		data = appendName(data, natives+uint32(i*8)+4, name)
	}
	for len(data)%cellBytes != 0 {
		data = append(data, 0)
	}
	cod := uint32(len(data))
	data = append(data, code...)
	dat := uint32(len(data))
	putHeader(data, uint32(len(data)))
	binary.LittleEndian.PutUint32(data[12:16], cod)
	binary.LittleEndian.PutUint32(data[16:20], dat)
	binary.LittleEndian.PutUint32(data[20:24], dat)
	binary.LittleEndian.PutUint32(data[24:28], dat+256)
	binary.LittleEndian.PutUint32(data[32:36], publics)
	binary.LittleEndian.PutUint32(data[36:40], natives)
	binary.LittleEndian.PutUint32(data[40:44], libs)
	binary.LittleEndian.PutUint32(data[44:48], libs)
	binary.LittleEndian.PutUint32(data[48:52], libs)
	binary.LittleEndian.PutUint32(data[52:56], libs)
	return data
}

func appendInstr(code []byte, op Opcode, params ...Cell) []byte {
	code = appendCell(code, Cell(op))
	for _, param := range params {
		code = appendCell(code, param)
	}
	return code
}

func appendCell(data []byte, cell Cell) []byte {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], uint32(int32(cell)))
	return append(data, buf[:]...)
}

func putHeader(data []byte, size uint32) {
	binary.LittleEndian.PutUint32(data[0:4], size)
	binary.LittleEndian.PutUint16(data[4:6], amxMagic32)
	data[6] = 8
	data[7] = 11
	binary.LittleEndian.PutUint16(data[10:12], 8)
	binary.LittleEndian.PutUint32(data[12:16], size)
	binary.LittleEndian.PutUint32(data[16:20], size)
	binary.LittleEndian.PutUint32(data[20:24], size)
	binary.LittleEndian.PutUint32(data[24:28], size+256)
	binary.LittleEndian.PutUint32(data[32:36], headerSize)
	binary.LittleEndian.PutUint32(data[36:40], headerSize+16)
	binary.LittleEndian.PutUint32(data[40:44], headerSize+24)
	binary.LittleEndian.PutUint32(data[44:48], headerSize+24)
	binary.LittleEndian.PutUint32(data[48:52], headerSize+24)
	binary.LittleEndian.PutUint32(data[52:56], headerSize+24)
}

func appendName(data []byte, stubNameOffset uint32, name string) []byte {
	off := uint32(len(data))
	binary.LittleEndian.PutUint32(data[stubNameOffset:stubNameOffset+4], off)
	data = append(data, name...)
	data = append(data, 0)
	binary.LittleEndian.PutUint32(data[0:4], uint32(len(data)))
	return data
}

package vm

import (
	"encoding/binary"
	"fmt"
)

const (
	amxMagic16     = 0xf1e2
	amxMagic32     = 0xf1e0
	amxMagic64     = 0xf1e1
	amxFlagCompact = 0x04
	amxFlagDebug   = 0x02
	headerSize     = 56
	minFileVersion = 6
	maxFileVersion = 9
	maxNameLength  = 31
)

type header struct {
	Size        uint32
	Magic       uint16
	FileVersion uint8
	AMXVersion  uint8
	Flags       uint16
	DefSize     uint16
	COD         uint32
	DAT         uint32
	HEA         uint32
	STP         uint32
	CIP         uint32
	Publics     uint32
	Natives     uint32
	Libraries   uint32
	PubVars     uint32
	Tags        uint32
	NameTable   uint32
}

func parseHeader(data []byte) (header, error) {
	if len(data) < headerSize {
		return header{}, fmt.Errorf("%w: file is shorter than AMX header", ErrInvalidAMX)
	}
	h := header{
		Size:        binary.LittleEndian.Uint32(data[0:4]),
		Magic:       binary.LittleEndian.Uint16(data[4:6]),
		FileVersion: data[6],
		AMXVersion:  data[7],
		Flags:       binary.LittleEndian.Uint16(data[8:10]),
		DefSize:     binary.LittleEndian.Uint16(data[10:12]),
		COD:         binary.LittleEndian.Uint32(data[12:16]),
		DAT:         binary.LittleEndian.Uint32(data[16:20]),
		HEA:         binary.LittleEndian.Uint32(data[20:24]),
		STP:         binary.LittleEndian.Uint32(data[24:28]),
		CIP:         binary.LittleEndian.Uint32(data[28:32]),
		Publics:     binary.LittleEndian.Uint32(data[32:36]),
		Natives:     binary.LittleEndian.Uint32(data[36:40]),
		Libraries:   binary.LittleEndian.Uint32(data[40:44]),
		PubVars:     binary.LittleEndian.Uint32(data[44:48]),
		Tags:        binary.LittleEndian.Uint32(data[48:52]),
		NameTable:   binary.LittleEndian.Uint32(data[52:56]),
	}
	if h.Magic == amxMagic16 || h.Magic == amxMagic64 {
		return header{}, fmt.Errorf("%w: magic 0x%04x", ErrUnsupportedCellSize, h.Magic)
	}
	if h.Magic != amxMagic32 {
		return header{}, fmt.Errorf("%w: bad magic 0x%04x", ErrInvalidAMX, h.Magic)
	}
	if !supportedVersion(h) {
		return header{}, fmt.Errorf("%w: unsupported file/amx version %d/%d", ErrInvalidAMX, h.FileVersion, h.AMXVersion)
	}
	if h.DefSize != 8 && h.DefSize != 24 {
		return header{}, fmt.Errorf("%w: unsupported definition size %d", ErrInvalidAMX, h.DefSize)
	}
	if h.Size == 0 || int(h.Size) > len(data) {
		return header{}, fmt.Errorf("%w: header size %d exceeds file size %d", ErrInvalidAMX, h.Size, len(data))
	}
	if h.COD < headerSize || h.COD > h.DAT || h.DAT > h.HEA || h.HEA > h.STP {
		return header{}, fmt.Errorf("%w: invalid segment ordering", ErrInvalidAMX)
	}
	if h.Publics > h.Natives || h.Natives > h.Libraries || h.Libraries > h.PubVars || h.PubVars > h.Tags {
		return header{}, fmt.Errorf("%w: invalid table ordering", ErrInvalidAMX)
	}
	if h.DefSize == 8 {
		if h.Tags > h.NameTable || h.NameTable > h.COD || int(h.NameTable)+2 > len(data) {
			return header{}, fmt.Errorf("%w: invalid name table", ErrInvalidAMX)
		}
		nameLength := binary.LittleEndian.Uint16(data[h.NameTable : h.NameTable+2])
		if nameLength > maxNameLength {
			return header{}, fmt.Errorf("%w: maximum name length %d exceeds %d", ErrInvalidAMX, nameLength, maxNameLength)
		}
	} else if h.Tags > h.COD {
		return header{}, fmt.Errorf("%w: legacy tag table exceeds code offset", ErrInvalidAMX)
	}
	return h, nil
}

func supportedVersion(h header) bool {
	fileVersionOK := h.FileVersion >= minFileVersion && h.FileVersion <= maxFileVersion
	amxVersionOK := h.AMXVersion >= minFileVersion
	return fileVersionOK && amxVersionOK
}

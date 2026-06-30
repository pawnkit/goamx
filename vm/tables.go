package vm

import (
	"encoding/binary"
	"fmt"
)

type funcStub struct {
	Address uint32
	Name    string
}

func parseFuncTable(data []byte, start, end uint32, defSize uint16) ([]funcStub, error) {
	if start > end || int(end) > len(data) {
		return nil, fmt.Errorf("%w: invalid function table bounds", ErrInvalidAMX)
	}
	if start == end {
		return nil, nil
	}
	if (end-start)%uint32(defSize) != 0 {
		return nil, fmt.Errorf("%w: function table is not aligned to definition size", ErrInvalidAMX)
	}
	count := int((end - start) / uint32(defSize))
	out := make([]funcStub, 0, count)
	for i := range count {
		off := int(start) + i*int(defSize)
		var name string
		if defSize == 8 {
			nameOff := binary.LittleEndian.Uint32(data[off+4 : off+8])
			var err error
			name, err = readCString(data, nameOff)
			if err != nil {
				return nil, err
			}
		} else {
			nameBytes := data[off+4 : off+int(defSize)]
			end := 0
			for end < len(nameBytes) && nameBytes[end] != 0 {
				end++
			}
			if end == len(nameBytes) {
				return nil, fmt.Errorf("%w: unterminated inline function name", ErrInvalidAMX)
			}
			name = string(nameBytes[:end])
		}
		out = append(out, funcStub{
			Address: binary.LittleEndian.Uint32(data[off : off+4]),
			Name:    name,
		})
	}
	return out, nil
}

func readCString(data []byte, off uint32) (string, error) {
	if int(off) >= len(data) {
		return "", fmt.Errorf("%w: string offset %d outside file", ErrInvalidAMX, off)
	}
	end := int(off)
	for end < len(data) && data[end] != 0 {
		end++
	}
	if end == len(data) {
		return "", fmt.Errorf("%w: unterminated string at offset %d", ErrInvalidAMX, off)
	}
	return string(data[off:uint32(end)]), nil
}

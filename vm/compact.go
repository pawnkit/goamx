package vm

import (
	"encoding/binary"
	"fmt"
)

func expandCompactImage(data []byte, h header) ([]byte, header, error) {
	if h.COD > h.Size || h.Size > uint32(len(data)) || h.HEA < h.COD {
		return nil, header{}, fmt.Errorf("%w: invalid compact segment bounds", ErrInvalidAMX)
	}
	if h.Size >= h.HEA {
		h.Flags &^= amxFlagCompact
		return data, h, nil
	}
	memSize := int(h.HEA - h.COD)
	if memSize%cellBytes != 0 {
		return nil, header{}, fmt.Errorf("%w: compact memory size %d is not cell-aligned", ErrInvalidAMX, memSize)
	}
	expandedBlock, err := expandCompactBlock(data[h.COD:h.Size], memSize)
	if err != nil {
		return nil, header{}, err
	}
	originalSize := h.Size
	debugData := append([]byte(nil), data[originalSize:]...)
	expanded := make([]byte, int(h.HEA), int(h.HEA)+len(debugData))
	copy(expanded, data[:h.COD])
	copy(expanded[h.COD:h.HEA], expandedBlock)
	expanded = append(expanded, debugData...)
	h.Size = h.HEA
	h.Flags &^= amxFlagCompact
	binary.LittleEndian.PutUint32(expanded[0:4], h.Size)
	binary.LittleEndian.PutUint16(expanded[8:10], h.Flags)
	return expanded, h, nil
}

func expandCompactBlock(compact []byte, memSize int) ([]byte, error) {
	out := make([]byte, memSize)
	codesize := len(compact)
	write := memSize
	for codesize > 0 {
		var value uint32
		shift := 0
		start := codesize
		for {
			codesize--
			if shift >= 32 {
				return nil, fmt.Errorf("%w: compact cell exceeds 32 bits", ErrInvalidAMX)
			}
			b := compact[codesize]
			value |= uint32(b&0x7f) << uint(shift)
			shift += 7
			if codesize == 0 || compact[codesize-1]&0x80 == 0 {
				break
			}
		}
		if compact[codesize]&0x40 != 0 && shift < 32 {
			value |= ^uint32(0) << uint(shift)
		}
		write -= cellBytes
		if write < 0 {
			return nil, fmt.Errorf("%w: compact output exceeds memory size", ErrInvalidAMX)
		}
		binary.LittleEndian.PutUint32(out[write:write+cellBytes], value)
		if start == codesize {
			return nil, fmt.Errorf("%w: compact decoder made no progress", ErrInvalidAMX)
		}
	}
	if write != 0 {
		return nil, fmt.Errorf("%w: compact output ended with %d bytes unwritten", ErrInvalidAMX, write)
	}
	return out, nil
}

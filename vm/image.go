package vm

type ImageInfo struct {
	Name           string
	FileVersion    uint8
	AMXVersion     uint8
	Flags          uint16
	DefinitionSize uint16
	CodeSize       int
	DataSize       int
	StackHeapSize  int
	NameLength     int
	Main           bool
}

func (vm *VM) Info() ImageInfo {
	nameLength := 0
	if vm.header.DefSize != 8 {
		nameLength = int(vm.header.DefSize) - cellBytes
	} else if int(vm.header.NameTable)+2 <= len(vm.image) {
		nameLength = int(vm.image[vm.header.NameTable]) | int(vm.image[vm.header.NameTable+1])<<8
	}
	return ImageInfo{
		Name:           vm.name,
		FileVersion:    vm.header.FileVersion,
		AMXVersion:     vm.header.AMXVersion,
		Flags:          vm.sourceFlags,
		DefinitionSize: vm.header.DefSize,
		CodeSize:       int(vm.header.DAT - vm.header.COD),
		DataSize:       int(vm.header.HEA - vm.header.DAT),
		StackHeapSize:  int(vm.header.STP - vm.header.HEA),
		NameLength:     nameLength,
		Main:           vm.header.CIP != ^uint32(0),
	}
}

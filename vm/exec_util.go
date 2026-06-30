package vm

func boolCell(ok bool) Cell {
	if ok {
		return 1
	}
	return 0
}

func compareBytes(left, right []byte) Cell {
	for i := 0; i < len(left) && i < len(right); i++ {
		if left[i] != right[i] {
			return Cell(int(left[i]) - int(right[i]))
		}
	}
	return Cell(len(left) - len(right))
}

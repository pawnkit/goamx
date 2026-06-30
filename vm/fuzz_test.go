package vm

import "testing"

func FuzzLoadBytesNeverPanics(f *testing.F) {
	f.Add([]byte("not an amx"))
	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = LoadBytes("fuzz.amx", data)
	})
}

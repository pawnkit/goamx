package amx

import "testing"

func TestInstrumentationDisabledByDefault(t *testing.T) {
	var runtime Runtime
	if err := runtime.emit(InstrumentationEvent{Kind: EventInstruction}); err != nil {
		t.Fatal(err)
	}
}

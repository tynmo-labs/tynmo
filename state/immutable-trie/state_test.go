package itrie

import (
	"testing"

	"tynmo/state"
)

func TestState(t *testing.T) {
	state.TestState(t, buildPreState)
}

func buildPreState(pre state.PreStates) state.Snapshot {
	storage := NewMemoryStorage()
	st := NewState(storage)
	snap := st.NewSnapshot()

	return snap
}

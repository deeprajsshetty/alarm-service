package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIsValid_ValidStates tests valid alarm states.
func TestIsValid_ValidStates(t *testing.T) {
	validStates := []AlarmState{Triggered, Active, ACKed, Cleared}

	for _, state := range validStates {
		assert.True(t, state.IsValid(), "Expected state %v to be valid", state)
	}
}

// TestIsValid_InvalidState tests invalid alarm states.
func TestIsValid_InvalidState(t *testing.T) {
	invalidState := AlarmState("InvalidState")
	assert.False(t, invalidState.IsValid(), "Expected invalid state to return false")
}
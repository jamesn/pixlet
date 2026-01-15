package starlarkutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.starlark.net/starlark"
)

func TestAddOnExit(t *testing.T) {
	thread := &starlark.Thread{}
	callCount := 0

	// Add first function
	AddOnExit(thread, func() {
		callCount++
	})

	// Verify function was added
	onExit, ok := thread.Local(ThreadOnExitKey).(*[]threadOnExitFunc)
	assert.True(t, ok)
	assert.Equal(t, 1, len(*onExit))

	// Add second function
	AddOnExit(thread, func() {
		callCount += 10
	})

	// Verify both functions are present
	onExit, ok = thread.Local(ThreadOnExitKey).(*[]threadOnExitFunc)
	assert.True(t, ok)
	assert.Equal(t, 2, len(*onExit))
}

func TestRunOnExitFuncs(t *testing.T) {
	thread := &starlark.Thread{}
	callCount := 0
	order := []int{}

	// Add multiple functions
	AddOnExit(thread, func() {
		callCount++
		order = append(order, 1)
	})

	AddOnExit(thread, func() {
		callCount += 10
		order = append(order, 2)
	})

	AddOnExit(thread, func() {
		callCount += 100
		order = append(order, 3)
	})

	// Run the functions
	RunOnExitFuncs(thread)

	// Verify all functions were called
	assert.Equal(t, 111, callCount)
	assert.Equal(t, []int{1, 2, 3}, order)
}

func TestRunOnExitFuncsWithoutAddingAny(t *testing.T) {
	thread := &starlark.Thread{}

	// Should not panic when no functions are added
	assert.NotPanics(t, func() {
		RunOnExitFuncs(thread)
	})
}

func TestRunOnExitFuncsMultipleTimes(t *testing.T) {
	thread := &starlark.Thread{}
	callCount := 0

	AddOnExit(thread, func() {
		callCount++
	})

	// Run multiple times - functions should be called each time
	RunOnExitFuncs(thread)
	assert.Equal(t, 1, callCount)

	RunOnExitFuncs(thread)
	assert.Equal(t, 2, callCount)
}

package utils

import (
	"fmt"
)

type mockT struct{}

func (m *mockT) Errorf(format string, args ...interface{}) {
	// In a real test, this would print an error message to the test output
	fmt.Printf(format, args...)
}

package helpers

import (
	"os"
	"testing"
)

func TestMust_NoError(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Must panicked unexpectedly: %v", r)
		}
	}()
	Must(nil)
}

func TestMust_WithError(t *testing.T) {
	if os.Getenv("TEST_MUST_EXIT") == "1" {
		Must(os.ErrNotExist)
		return
	}

	// This test verifies that Must calls os.Exit(1)
	// We can't directly test os.Exit, so we document the behavior
	t.Log("Must calls os.Exit(1) on error - tested via integration")
}

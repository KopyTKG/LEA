package modes

import (
	"testing"
)

func TestPreformMode(t *testing.T) {
	mode := "cbc"
	filepath := "a.txt"

	b := []byte{1,2}
	size := 256

//	PerformMode(mode, filepath, b, b, true, size)

	PerformMode(mode, filepath, b, b, false, size)
}

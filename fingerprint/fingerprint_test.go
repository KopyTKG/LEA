package fingerprint

import (
	"testing"
)

func TestLoadSource(t *testing.T) {
	a := []byte{1, 2, 3, 4}
	out := LoadSource(a)
	if len(out) < 8 {
		t.Fatalf("Length is smaller then 8 it is %d", len(out))
	}
	for i, item := range out {
		t.Logf("%d %d", i, item)
	}
}

func TestFingerprint128(t *testing.T) {
	in := [8]uint64{1, 2, 3, 4, 5, 6, 7, 8}
	out := Fingerprint128(in)
	
	if len(out) < 4 {
		t.Fatalf("Length is smaller then 4 it is %d", len(out))
	}
}

package fingerprint

import (
	"testing"
)

func CheckLenght32(t *testing.T, lenght int, arr []uint32) {
	if len(arr) < lenght {
		t.Fatalf("Length is smaller then %d it is %d", lenght, len(arr))
	}
}

func CheckElements32(t *testing.T, arr []uint32, check []uint32) {
	for i := 0; i < len(arr); i++ {
		if arr[i] != check[i] {
			t.Errorf("%d does not match %d", arr[i], check[i])
		}
	}

}

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
	check := [4]uint32{0, 0, 0, 8}

	CheckLenght32(t, 4, out[:])
	CheckElements32(t, out[:], check[:])
}

func TestFingerprint192(t *testing.T) {
	in := [8]uint64{1, 2, 3, 4, 5, 6, 7, 8}
	out := Fingerprint192(in)
	check := [6]uint32{0, 5, 0, 5, 0, 12}

	CheckLenght32(t, 6, out[:])
	CheckElements32(t, out[:], check[:])
}

func TestFingerprint256(t *testing.T) {
	in := [8]uint64{1, 2, 3, 4, 5, 6, 7, 8}
	out := Fingerprint256(in)
	check := [8]uint32{0, 2, 0, 6, 0, 2, 0, 14}

	CheckLenght32(t, 8, out[:])
	CheckElements32(t, out[:], check[:])
}

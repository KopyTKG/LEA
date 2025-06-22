package schedule

import (
	"fmt"
	"lea/bitops"
)

const (
	// Chunk CHUNK256s
	CHUNK128 = 4 // 32bit x 4
	CHUNK192 = 6 // 32bit x 6
	CHUNK256 = 8 // 32bit x 8

	// Rounds
	ROUNDS128 = 24
	ROUNDS192 = 28
	ROUNDS256 = 32

	// Rk lengths
	KEY128 = 144 // 24 rounds x 6 words
	KEY192 = 168 // 28 rounds x 6 words
	KEY256 = 192 // 32 rounds x 6 words
)

type Chunk128 [CHUNK128]uint32
type Chunk192 [CHUNK192]uint32
type Chunk256 [CHUNK256]uint32

type Rk128 [KEY128]uint32
type Rk192 [KEY192]uint32
type Rk256 [KEY256]uint32

func KeySchedule(size int, key, seed []uint32) ([]uint32, error) {
	switch size {
	case 128:
		k := Chunk128{}
		s := Chunk128{}
		copy(k[:CHUNK128], key)
		copy(s[:CHUNK128], seed)
		arr, err := Schedule128(k, s)
		return arr[:], err
	case 192:
		k := Chunk192{}
		s := Chunk192{}
		copy(k[:CHUNK192], key)
		copy(s[:CHUNK192], seed)
		arr, err := Schedule192(k, s)
		return arr[:], err
	case 256:
		k := Chunk256{}
		s := Chunk256{}
		copy(k[:CHUNK256], key)
		copy(s[:CHUNK256], seed)
		arr, err := Schedule256(k, s)
		return arr[:], err
	default:
		return []uint32{}, fmt.Errorf("invalid key size %d, must be 128, 192 or 256", size)
	}
}

// Function failsafe
func checkLen(item []uint32, lenght int) (bool, error) {
	if len(item) == lenght {
		return true, nil
	}
	return false, fmt.Errorf("key or seed length is not valid, expected %d, got %d", lenght, len(item))
}

func Schedule128(key, seed Chunk128) (*Rk128, error) {
	// Lenght validation
	_, err := checkLen(key[:], CHUNK128)
	if err != nil {
		return nil, err
	}

	_, err = checkLen(seed[:], CHUNK128)
	if err != nil {
		return nil, err
	}

	rk := Rk128{}
	for i := range CHUNK128 {
		rk[i] = key[i]
	}

	// Generate round key
	var rkT = make([]uint32, CHUNK128)
	copy(rkT, key[:CHUNK128])
	for i := range ROUNDS128 {
		t0 := bitops.ShiftLeft32(seed[i%CHUNK128], uint(i))
		t1 := bitops.ShiftLeft32(t0, 1)
		t2 := bitops.ShiftLeft32(t1, 1)
		t3 := bitops.ShiftLeft32(t2, 1)

		rkT[0] = bitops.RotateLeft32(rkT[0]+t0, 1)
		rkT[1] = bitops.RotateLeft32(rkT[1]+t1, 3)
		rkT[2] = bitops.RotateLeft32(rkT[2]+t2, 6)
		rkT[3] = bitops.RotateLeft32(rkT[3]+t3, 11)

		j := 6 * i
		rk[j+0] = rkT[0]
		rk[j+1] = rkT[1]
		rk[j+2] = rkT[2]
		rk[j+3] = rkT[1]
		rk[j+4] = rkT[3]
		rk[j+5] = rkT[1]
	}
	return &rk, nil
}

func Schedule192(key, seed Chunk192) (*Rk192, error) {
	// Lenght validation
	_, err := checkLen(key[:], CHUNK192)
	if err != nil {
		return nil, err
	}

	_, err = checkLen(seed[:], CHUNK192)
	if err != nil {
		return nil, err
	}

	rk := Rk192{}
	for i := range CHUNK192 {
		rk[i] = key[i]
	}

	// Generate round key
	var rkT = make([]uint32, CHUNK192)
	copy(rkT, key[:CHUNK192])
	for i := range ROUNDS192 {
		s := i % 32
		t0 := bitops.ShiftLeft32(seed[i%CHUNK192], uint(s))
		t1 := bitops.ShiftLeft32(t0, 1)
		t2 := bitops.ShiftLeft32(t0, 2)
		t3 := bitops.ShiftLeft32(t0, 3)
		t4 := bitops.ShiftLeft32(t0, 4)
		t5 := bitops.ShiftLeft32(t0, 5)

		j := 6 * i
		rk[j+0] = bitops.RotateLeft32(rkT[0]+t0, 1)
		rk[j+1] = bitops.RotateLeft32(rkT[1]+t1, 3)
		rk[j+2] = bitops.RotateLeft32(rkT[2]+t2, 6)
		rk[j+3] = bitops.RotateLeft32(rkT[3]+t3, 11)
		rk[j+4] = bitops.RotateLeft32(rkT[4]+t4, 13)
		rk[j+5] = bitops.RotateLeft32(rkT[5]+t5, 17)
	}
	return &rk, nil
}

func Schedule256(key, seed Chunk256) (*Rk256, error) {
	// Lenght validation
	_, err := checkLen(key[:], CHUNK256)
	if err != nil {
		return nil, err
	}

	_, err = checkLen(seed[:], CHUNK256)
	if err != nil {
		return nil, err
	}

	rk := Rk256{}
	for i := range CHUNK256 {
		rk[i] = key[i]
	}

	// Generate round key
	var rkT = make([]uint32, CHUNK256)
	copy(rkT, key[:CHUNK256])
	for i := range ROUNDS256 {
		s := i % 32
		t0 := bitops.ShiftLeft32(seed[i%CHUNK256], uint(s))
		t1 := bitops.ShiftLeft32(t0, 1)
		t2 := bitops.ShiftLeft32(t0, 2)
		t3 := bitops.ShiftLeft32(t0, 3)
		t4 := bitops.ShiftLeft32(t0, 4)
		t5 := bitops.ShiftLeft32(t0, 5)

		j := 6 * i
		rk[j+0] = bitops.RotateLeft32(rkT[(j+0)%CHUNK256]+t0, 1)
		rk[j+1] = bitops.RotateLeft32(rkT[(j+1)%CHUNK256]+t1, 3)
		rk[j+2] = bitops.RotateLeft32(rkT[(j+2)%CHUNK256]+t2, 6)
		rk[j+3] = bitops.RotateLeft32(rkT[(j+3)%CHUNK256]+t3, 11)
		rk[j+4] = bitops.RotateLeft32(rkT[(j+4)%CHUNK256]+t4, 13)
		rk[j+5] = bitops.RotateLeft32(rkT[(j+5)%CHUNK256]+t5, 17)
	}
	return &rk, nil
}

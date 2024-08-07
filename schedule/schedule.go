package schedule

import (
	"lea/bitops"
)

func KeySchedule(size int, key, seed []uint32) []uint32 {
	switch size {
	case 128:
		k := [4]uint32{}
		s := [4]uint32{}
		copy(k[:4], key)
		copy(s[:4], seed)
		arr := Schedule128(k, s)
		return arr[:]
	case 192:
		k := [6]uint32{}
		s := [6]uint32{}
		copy(k[:6], key)
		copy(s[:6], seed)
		arr := Schedule192(k, s)
		return arr[:]
	case 256:
		k := [8]uint32{}
		s := [8]uint32{}
		copy(k[:8], key)
		copy(s[:8], seed)
		arr := Schedule256(k, s)
		return arr[:]
	default:
		return []uint32{}
	}
}

// Function failsafe
func checkLen(item []uint32, lenght int) bool {
	if len(item) == lenght {
		return true
	}
	panic("Not valid lenght")
}

func Schedule128(key, seed [4]uint32) [144]uint32 {
	size := 4
	// Lenght validation
	checkLen(key[:], size)
	checkLen(seed[:], size)

	rk := [144]uint32{}
	for i := 0; i < size; i++ {
		rk[i] = key[i]
	}

	// Generate round key
	var rkT = make([]uint32, size)
	copy(rkT, key[:size])
	for i := 0; i < 24; i++ {
		t0 := bitops.ShiftLeft32(seed[i%size], uint(i))
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
	return rk
}

func Schedule192(key, seed [6]uint32) [168]uint32 {
	size := 6
	// Lenght validation
	checkLen(key[:], size)
	checkLen(seed[:], size)

	rk := [168]uint32{}
	for i := 0; i < size; i++ {
		rk[i] = key[i]
	}

	// Generate round key
	var rkT = make([]uint32, size)
	copy(rkT, key[:size])
	for i := 0; i < 28; i++ {
		s := i % 32
		t0 := bitops.ShiftLeft32(seed[i%size], uint(s))
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
	return rk
}

func Schedule256(key, seed [8]uint32) [192]uint32 {
	size := 8
	// Lenght validation
	checkLen(key[:], size)
	checkLen(seed[:], size)

	rk := [192]uint32{}
	for i := 0; i < size; i++ {
		rk[i] = key[i]
	}

	// Generate round key
	var rkT = make([]uint32, size)
	copy(rkT, key[:size])
	for i := 0; i < 32; i++ {
		s := i % 32
		t0 := bitops.ShiftLeft32(seed[i%size], uint(s))
		t1 := bitops.ShiftLeft32(t0, 1)
		t2 := bitops.ShiftLeft32(t0, 2)
		t3 := bitops.ShiftLeft32(t0, 3)
		t4 := bitops.ShiftLeft32(t0, 4)
		t5 := bitops.ShiftLeft32(t0, 5)

		j := 6 * i
		rk[j+0] = bitops.RotateLeft32(rkT[(j+0)%size]+t0, 1)
		rk[j+1] = bitops.RotateLeft32(rkT[(j+1)%size]+t1, 3)
		rk[j+2] = bitops.RotateLeft32(rkT[(j+2)%size]+t2, 6)
		rk[j+3] = bitops.RotateLeft32(rkT[(j+3)%size]+t3, 11)
		rk[j+4] = bitops.RotateLeft32(rkT[(j+4)%size]+t4, 13)
		rk[j+5] = bitops.RotateLeft32(rkT[(j+5)%size]+t5, 17)
	}
	return rk
}

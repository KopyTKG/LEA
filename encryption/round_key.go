package encryption

import "lea/bitops"

var constval = [8]uint32{0xC3EFE9DB, 0x44626B02, 0x79E27C8A, 0x78DF30EC, 0x715EA49E, 0xC785DA0A, 0xE04EF22A, 0xE5C40957}

func Generate(key [16]uint32) [144]uint32 {
	if len(key) != 16 {
		panic("Key must be an array of 16 uint32s")
	}

	// Initialize round key array
	rk := [144]uint32{}

	// Copy key to round key
	for i := 0; i < 4; i++ {
		rk[i] = key[i]
	}

	// Generate round key
	var rkT = make([]uint32, 4)
	copy(rkT, key[:4]) // Copy first four elements to temporary array

	for i := 0; i < 24; i++ {
		t0 := bitops.RotateLeft32(constval[i%4], uint(i))
		t1 := bitops.RotateLeft32(t0, 1)
		t2 := bitops.RotateLeft32(t1, 1)
		t3 := bitops.RotateLeft32(t2, 1)

		rkT[0] = bitops.RotateLeft32(rkT[0]+t0, 1)
		rkT[1] = bitops.RotateLeft32(rkT[1]+t1, 3)
		rkT[2] = bitops.RotateLeft32(rkT[2]+t2, 6)
		rkT[3] = bitops.RotateLeft32(rkT[3]+t3, 11)

		j := 6 * i
		rk[j+0] = rkT[0]
		rk[j+1] = rkT[1]
		rk[j+2] = rkT[2]
		rk[j+3] = rkT[1] // Assuming the repeated use of rkT[1] is intentional
		rk[j+4] = rkT[3]
		rk[j+5] = rkT[1] // Repeated again
	}

	return rk
}

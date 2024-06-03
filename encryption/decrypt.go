package encryption

import (
	"lea/bitops"
)

func DecryptBlock(block [4]uint32, rk [144]uint32) [4]uint32 {
	var enc_block = [4]uint32{}
	for i := 0; i < 4; i++ {
		enc_block[i] = block[i]
	}
	for i := 23; i >= 0; i-- {
		rkI := 6 * i
		b0 := enc_block[3]
		b1 := bitops.WrappedSub32(bitops.RotateRight32(enc_block[0], 9), (b0^rk[rkI])) ^ rk[rkI+1]
		b2 := bitops.WrappedSub32(bitops.RotateLeft32(enc_block[1], 5), (b1^rk[rkI+2])) ^ rk[rkI+3]
		b3 := bitops.WrappedSub32(bitops.RotateLeft32(enc_block[2], 3), (b2^rk[rkI+4])) ^ rk[rkI+5]

		enc_block[0], enc_block[1], enc_block[2], enc_block[3] = b0, b1, b2, b3
	}

	return enc_block
}

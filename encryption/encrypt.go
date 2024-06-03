package encryption

import (
	"lea/bitops"
)

func EncryptBlock(block [4]uint32, rk [144]uint32) [4]uint32 {
	var enc_block = [4]uint32{}
	for i := 0; i < 4; i++ {
		enc_block[i] = block[i]
	}
	for i := 0; i < 24; i++ {
		rkI := 6 * i
		b0 := bitops.RotateLeft32(bitops.WrappedAdd32((enc_block[0]^rk[rkI]), (enc_block[1]^rk[rkI+1])), 9)
		b1 := bitops.RotateRight32(bitops.WrappedAdd32((enc_block[1]^rk[rkI+2]), (enc_block[2]^rk[rkI+3])), 5)
		b2 := bitops.RotateRight32(bitops.WrappedAdd32((enc_block[2]^rk[rkI+4]), (enc_block[3]^rk[rkI+5])), 3)
		b3 := enc_block[0]

		enc_block[0], enc_block[1], enc_block[2], enc_block[3] = b0, b1, b2, b3
	}

	return enc_block
}

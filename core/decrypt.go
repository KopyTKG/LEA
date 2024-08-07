package core

import (
	"lea/bitops"
	"lea/types"
)

func SelectDecrypt(block [4]uint32, rk []uint32, size int) []uint32 {
	switch size {
		case 128:
			b := Decrypt128(block, types.Rk128(rk))
			return b[:]
		case 192:
			b := Decrypt192(block, types.Rk192(rk))
			return b[:]
		case 256:
			b := Decrypt256(block, types.Rk256(rk))
			return b[:]
		default:
			return []uint32{}
	}
}

func decRound(block *[4]uint32, rk []uint32, i int) {
	rkI := 6 * i
	b0 := block[3]
	b1 := bitops.WrappedSub32(bitops.RotateRight32(block[0], 9), (b0^rk[rkI])) ^ rk[rkI+1]
	b2 := bitops.WrappedSub32(bitops.RotateLeft32(block[1], 5), (b1^rk[rkI+2])) ^ rk[rkI+3]
	b3 := bitops.WrappedSub32(bitops.RotateLeft32(block[2], 3), (b2^rk[rkI+4])) ^ rk[rkI+5]

	block[0], block[1], block[2], block[3] = b0, b1, b2, b3
}

func Decrypt128(block [4]uint32, rk types.Rk128) [4]uint32 {
	for i := 23; i >= 0; i-- {
		decRound(&block, rk[:], i)
	}
	return block
}

func Decrypt192(block [4]uint32, rk types.Rk192) [4]uint32 {
	for i := 27; i >= 0; i-- {
		decRound(&block, rk[:], i)
	}
	return block
}

func Decrypt256(block [4]uint32, rk types.Rk256) [4]uint32 {
	for i := 31; i >= 0; i-- {
		decRound(&block, rk[:], i)
	}
	return block
}

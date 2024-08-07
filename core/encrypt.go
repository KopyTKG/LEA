package core

import (
	"lea/bitops"
	"lea/types"
)

func SelectEncrypt(block [4]uint32, rk []uint32, size int) []uint32 {
	switch size {
		case 128:
			b := Encrypt128(block, types.Rk128(rk))
			return b[:]
		case 192:
			b := Encrypt192(block, types.Rk192(rk))
			return b[:]
		case 256:
			b := Encrypt256(block, types.Rk256(rk))
			return b[:]
		default:
			return []uint32{}
	}
}


func encRound(block *[4]uint32, rk []uint32, i int) {
	rkI := 6 * i
	b0 := bitops.RotateLeft32(bitops.WrappedAdd32((block[0]^rk[rkI]), (block[1]^rk[rkI+1])), 9)
	b1 := bitops.RotateRight32(bitops.WrappedAdd32((block[1]^rk[rkI+2]), (block[2]^rk[rkI+3])), 5)
	b2 := bitops.RotateRight32(bitops.WrappedAdd32((block[2]^rk[rkI+4]), (block[3]^rk[rkI+5])), 3)
	b3 := block[0]

	block[0], block[1], block[2], block[3] = b0, b1, b2, b3

}

func Encrypt128(block [4]uint32, rk types.Rk128) [4]uint32 {
	for i := 0; i < 24; i++ {
		encRound(&block, rk[:], i)
	}
	return block
}

func Encrypt192(block [4]uint32, rk types.Rk192) [4]uint32 {
	for i := 0; i < 28; i++ {
		encRound(&block, rk[:], i)
	}

	return block
}

func Encrypt256(block [4]uint32, rk types.Rk256) [4]uint32 {
	for i := 0; i < 32; i++ {
		encRound(&block, rk[:], i)
	}

	return block
}

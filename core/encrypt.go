package core

import (
	"lea/bitops"
	"lea/schedule"
	"lea/types"
)

func SelectEncrypt(block [4]uint32, rk []uint32, size int) [4]uint32 {
	switch size {
	case 128:
		return Encrypt128(block, types.Rk128(rk))
	case 192:
		return Encrypt192(block, types.Rk192(rk))
	case 256:
		return Encrypt256(block, types.Rk256(rk))
	default:
		return [4]uint32{}
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
	for i := range schedule.ROUNDS128 {
		encRound(&block, rk[:], i)
	}
	return block
}

func Encrypt192(block [4]uint32, rk types.Rk192) [4]uint32 {
	for i := range schedule.ROUNDS192 {
		encRound(&block, rk[:], i)
	}

	return block
}

func Encrypt256(block [4]uint32, rk types.Rk256) [4]uint32 {
	for i := range schedule.ROUNDS256 {
		encRound(&block, rk[:], i)
	}

	return block
}

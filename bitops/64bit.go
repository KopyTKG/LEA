package bitops

/* Bitwise operations for 64bit */

func WrappedSub64(x, y uint64) uint64 {
	return x - y
}

func WrappedAdd64(x, y uint64) uint64 {
	return x + y
}

func RotateLeft64(x uint64, n uint) uint64 {
	return (x << n) | (x >> (64 - n))
}

func RotateRight64(x uint64, n uint) uint64 {
	return (x >> n) | (x << (64 - n))
}

func Reverse64(x uint64) uint64 {
	x = (x&0x5555555555555555)<<1 | (x&0xAAAAAAAAAAAAAAAA)>>1
	x = (x&0x3333333333333333)<<2 | (x&0xCCCCCCCCCCCCCCCC)>>2
	x = (x&0x0F0F0F0F0F0F0F0F)<<4 | (x&0xF0F0F0F0F0F0F0F0)>>4
	x = (x&0x00FF00FF00FF00FF)<<8 | (x&0xFF00FF00FF00FF00)>>8
	x = (x&0x0000FFFF0000FFFF)<<16 | (x&0xFFFF0000FFFF0000)>>16
	x = (x&0x00000000FFFFFFFF)<<32 | (x&0xFFFFFFFF00000000)>>32
	return x
}

func CountBits64(x uint64) int {
	x = x - ((x >> 1) & 0x5555555555555555)
	x = (x & 0x3333333333333333) + ((x >> 2) & 0x3333333333333333)
	x = (x + (x >> 4)) & 0x0F0F0F0F0F0F0F0F
	x = x + (x >> 8)
	x = x + (x >> 16)
	x = x + (x >> 32)
	return int(x & 0x7F)
}

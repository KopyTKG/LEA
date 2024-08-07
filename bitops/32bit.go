package bitops

/* Bitwise operations for 32bit */

func WrappedSub32(x, y uint32) uint32 {
	return x - y
}

func WrappedAdd32(x, y uint32) uint32 {
	return x + y
}

func RotateLeft32(x uint32, n uint) uint32 {
	return (x << n) | (x >> (32 - n))
}

func RotateRight32(x uint32, n uint) uint32 {
	return (x >> n) | (x << (32 - n))
}

func Reverse32(x uint32) uint32 {
	x = (x&0x55555555)<<1 | (x&0xAAAAAAAA)>>1
	x = (x&0x33333333)<<2 | (x&0xCCCCCCCC)>>2
	x = (x&0x0F0F0F0F)<<4 | (x&0xF0F0F0F0)>>4
	x = (x&0x00FF00FF)<<8 | (x&0xFF00FF00)>>8
	x = (x&0x0000FFFF)<<16 | (x&0xFFFF0000)>>16
	return x
}

func CountBits32(x uint32) int {
	x = x - ((x >> 1) & 0x55555555)
	x = (x & 0x33333333) + ((x >> 2) & 0x33333333)
	x = (x + (x >> 4)) & 0x0F0F0F0F
	x = x + (x >> 8)
	x = x + (x >> 16)
	return int(x & 0x3F)
}

func ShiftLeft32(x uint32, n uint) uint32 {
	return x << uint32(n) 
}

func ShiftRight32(x uint32, n uint) uint32 {
	return x >> uint32(n)
}

func MultiXOR32(a, b [4]uint32) [4]uint32 {
	var result [4]uint32
	for i := 0; i < len(a); i++ {
		result[i] = a[i] ^ b[i]
	}
	return result
}

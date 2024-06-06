package bitops

/* Bitwise operations for 16bit */

func WrappedSub16(x, y uint16) uint16 {
	return x - y
}

func WrappedAdd16(x, y uint16) uint16 {
	return x + y
}

func RotateLeft16(x uint16, n uint) uint16 {
	return (x << n) | (x >> (16 - n))
}

func RotateRight16(x uint16, n uint) uint16 {
	return (x >> n) | (x << (16 - n))
}

func Reverse16(x uint16) uint16 {
	x = (x&0x5555)<<1 | (x&0xAAAA)>>1
	x = (x&0x3333)<<2 | (x&0xCCCC)>>2
	x = (x&0x0F0F)<<4 | (x&0xF0F0)>>4
	x = (x&0x00FF)<<8 | (x&0xFF00)>>8
	return x
}

func CountBits16(x uint16) int {
	x = x - ((x >> 1) & 0x5555)
	x = (x & 0x3333) + ((x >> 2) & 0x3333)
	x = (x + (x >> 4)) & 0x0F0F
	x = x + (x >> 8)
	return int(x & 0x1F)
}


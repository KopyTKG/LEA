package bitops

/* Bitwise operations for 8bit */

func WrappedSub8(x, y uint8) uint8 {
	return x - y
}

func WrappedAdd8(x, y uint8) uint8 {
	return x + y
}

func RotateLeft8(x uint8, n uint) uint8 {
	return (x << n) | (x >> (8 - n))
}

func RotateRight8(x uint8, n uint) uint8 {
	return (x >> n) | (x << (8 - n))
}

func Reverse8(x uint8) uint8 {
	x = (x&0x55)<<1 | (x&0xAA)>>1
	x = (x&0x33)<<2 | (x&0xCC)>>2
	x = (x&0x0F)<<4 | (x&0xF0)>>4
	return x
}


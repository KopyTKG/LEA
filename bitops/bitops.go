package bitops

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

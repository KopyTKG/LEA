package fingerprint

import (
	"encoding/binary"
	"golang.org/x/crypto/sha3"
)

var UPPERMASK uint64 = 0xFFFFFFFF00000000
var LOWERMASK uint64 = 0xFFFFFFFF

func LoadSource(data []byte) [8]uint64 {
	hasher := sha3.New512()
	_, err := hasher.Write(data)

	if err != nil {
		panic(err)
	}

	sum := hasher.Sum(nil)

	var hashArray [8]uint64
	for i := 0; i < 8; i++ {
		hashArray[i] = binary.BigEndian.Uint64(sum[i*8 : (i+1)*8])
	}
	return hashArray
}

func Fingerprint128(source [8]uint64) [4]uint32 {
	base := [4]uint32{}

	left := [2]uint64{source[0] ^ source[2], source[1] ^ source[3]}
	right := [2]uint64{source[4] ^ source[6], source[5] ^ source[7]}

	lowest := [2]uint64{left[0] ^ right[0], left[1] ^ right[1]}

	base[0] = uint32((lowest[0] & UPPERMASK) >> 32)
	base[1] = uint32((lowest[0] & LOWERMASK))
	base[2] = uint32((lowest[1] & UPPERMASK) >> 32)
	base[3] = uint32((lowest[1] & LOWERMASK))

	return base
}

func Fingerprint192(source [8]uint64) [6]uint32 {
	base := [6]uint32{}
	left := source[2] ^ source[4]
	right := source[3] ^ source[5]

	inner := [6]uint64{}

	inner[0] = source[0] ^ left
	inner[1] = source[1]
	inner[2] = left
	inner[3] = right
	inner[4] = source[6]
	inner[5] = source[7] ^ right

	base[0] = uint32(((inner[0] ^ inner[3]) & UPPERMASK) >> 32)
	base[1] = uint32((inner[0] ^ inner[3]) & LOWERMASK)
	base[2] = uint32(((inner[1] ^ inner[4]) & UPPERMASK) >> 32)
	base[3] = uint32((inner[1] ^ inner[4]) & LOWERMASK)
	base[4] = uint32(((inner[2] ^ inner[5]) & UPPERMASK) >> 32)
	base[5] = uint32((inner[2] ^ inner[5]) & LOWERMASK)

	return base
}

func Fingerprint256(source [8]uint64) [8]uint32 {
	base := [8]uint32{}

	left := [2]uint64{source[0] ^ source[2], source[1] ^ source[3]}
	right := [2]uint64{source[4] ^ source[6], source[5] ^ source[7]}

	base[0] = uint32(((left[0]) & UPPERMASK) >> 32)
	base[1] = uint32((left[0]) & LOWERMASK)
	base[2] = uint32(((left[1]) & UPPERMASK) >> 32)
	base[3] = uint32((left[1]) & LOWERMASK)
	base[4] = uint32(((right[0]) & UPPERMASK) >> 32)
	base[5] = uint32((right[0]) & LOWERMASK)
	base[6] = uint32(((right[1]) & UPPERMASK) >> 32)
	base[7] = uint32((right[1]) & LOWERMASK)
	return base
}

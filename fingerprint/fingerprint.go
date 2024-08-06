package fingerprint

import (
	"encoding/binary"
	"golang.org/x/crypto/sha3"
)

var UPPERMASK uint = 0xFFFFFFFF00000000
var LOWERMASK uint = 0x00000000FFFFFFFF


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

	base[0] = uint32((lowest[0] & 0xFFFFFFFF00000000) >> 32)
	base[1] = uint32((lowest[0] & 0x00000000FFFFFFFF))
	base[2] = uint32((lowest[1] & 0xFFFFFFFF00000000) >> 32)
	base[3] = uint32((lowest[1] & 0x00000000FFFFFFFF))

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

	base[0] = uint32(((inner[0] ^ inner[3]) & 0xFFFFFFFF00000000) >> 32)
	base[1] = uint32((inner[0] ^ inner[3]) & 0x00000000FFFFFFFF)
	base[2] = uint32(((inner[1] ^ inner[4]) & 0xFFFFFFFF00000000) >> 32)
	base[3] = uint32((inner[1] ^ inner[4]) & 0x00000000FFFFFFFF)
	base[4] = uint32(((inner[2] ^ inner[5]) & 0xFFFFFFFF00000000) >> 32)
	base[5] = uint32((inner[2] ^ inner[5]) & 0x00000000FFFFFFFF)

	return base
}

func Fingerprint256(source [8]uint64) [8]uint32 {
	base := [8]uint32{}

	left := [2]uint64{source[0] ^ source[2], source[1] ^ source[3]}
	right := [2]uint64{source[4] ^ source[6], source[5] ^ source[7]}

	base[0] = uint32(((left[0]) & 0xFFFFFFFF00000000) >> 32)
	base[1] = uint32((left[0]) & 0x00000000FFFFFFFF)
	base[2] = uint32(((left[1]) & 0xFFFFFFFF00000000) >> 32)
	base[3] = uint32((left[1]) & 0x00000000FFFFFFFF)
	base[4] = uint32(((right[0]) & 0xFFFFFFFF00000000) >> 32)
	base[5] = uint32((right[0]) & 0x00000000FFFFFFFF)
	base[6] = uint32(((right[1]) & 0xFFFFFFFF00000000) >> 32)
	base[7] = uint32((right[1]) & 0x00000000FFFFFFFF)
	return base
}

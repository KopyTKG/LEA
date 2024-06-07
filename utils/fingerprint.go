package utils

import (
	"encoding/binary"
	"golang.org/x/crypto/sha3"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// FingerPrint512 returns the SHA3-512 hash of the input data as an array of uint32.
func FingerPrint512(data []byte) [16]uint32 {
    hasher := sha3.New512()
    _, err := hasher.Write(data)
    check(err)
    sum := hasher.Sum(nil)

    var hashArray [16]uint32
    for i := 0; i < 16; i++ {
        hashArray[i] = binary.BigEndian.Uint32(sum[i*4 : (i+1)*4])
    }
    return hashArray
}

// FingerPrint256 returns the SHA3-256 hash of the input data as an array of uint32.
func FingerPrint256(data []byte) [8]uint32 {
    hasher := sha3.New256()
    _, err := hasher.Write(data)
    check(err)
    sum := hasher.Sum(nil)

    var hashArray [8]uint32
    for i := 0; i < 8; i++ {
        hashArray[i] = binary.BigEndian.Uint32(sum[i*4 : (i+1)*4])
    }
    return hashArray
}

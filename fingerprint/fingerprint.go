package fingerprint

import (
	//"lea/bitops"
	"encoding/binary"
	"golang.org/x/crypto/sha3"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func LoadSource(data []byte) [8]uint64 {
    hasher := sha3.New512()
    _, err := hasher.Write(data)
    check(err)
    sum := hasher.Sum(nil)

    var hashArray [8]uint64
    for i := 0; i < 8; i++ {
        hashArray[i] = binary.BigEndian.Uint64(sum[i*8 : (i+1)*8])
    }
    return hashArray
}

func Fingerprint128(source [8]uint64) [4]uint32 {
 base := [4]uint32{}
 return base
}

func Fingerprint192(source [8]uint64) [6]uint32 {
 base := [6]uint32{}
 return base
}

func Fingerprint256(source [8]uint64) [8]uint32{
 base := [8]uint32{}
 return base
}

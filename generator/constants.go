package generator

import (
	"lea/stream"
	"math/rand"
	"time"
)



func GenerateConstants() {
	rand.Seed(time.Now().UnixNano())
	var seed = []uint32{}

	for i := 0; i < 8; i++ {
		seed = append(seed, uint32(rand.Intn(0xFFFFFFFF)))	
	}

	stream.WriteBinaryStream("/tmp/seed", seed)

}

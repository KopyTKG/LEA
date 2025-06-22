package modes

import (
	"lea/bitops"
	"lea/core"
	"lea/stream"

	"github.com/kopytkg/golog"
)

func encryptCFB(filePath string, prev *[4]uint32, keySegments []uint32, chunks [4]uint32, size int) {
	encB := [4]uint32(core.SelectEncrypt(*prev, keySegments, size))

	*prev = bitops.MultiXOR32(encB, chunks)

	if err := stream.WriteBinaryStream(filePath, *prev); err != nil {
		golog.Errorf("Error writing to binary stream: %v\n", err)
	}
}

func decryptCFB(filePath string, prev *[4]uint32, keySegments []uint32, chunks [4]uint32, size int) {
	encB := [4]uint32(core.SelectEncrypt(*prev, keySegments, size))

	text := bitops.MultiXOR32(encB, chunks)

	if err := stream.WriteBinaryStream(filePath, text); err != nil {
		golog.Errorf("Error prepending to binary stream: %v\n", err)
	}

	*prev = chunks

}

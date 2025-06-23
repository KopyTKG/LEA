package modes

import (
	"lea/bitops"
	"lea/core"
	"lea/stream"

	"github.com/kopytkg/golog"
)

func encryptCTR(filePath string, keySegments []uint32, chunks [4]uint32, size int, counter *[4]uint32) {
	encB := [4]uint32(core.SelectEncrypt(*counter, keySegments, size))

	prev := bitops.MultiXOR32(encB, chunks)

	if err := stream.WriteBinaryStream(filePath, prev); err != nil {
		golog.Errorf("Error writing to binary stream: %v\n", err)
	}
}

func decryptCTR(filePath string, keySegments []uint32, chunks [4]uint32, size int, counter *[4]uint32) {
	encB := [4]uint32(core.SelectEncrypt(*counter, keySegments, size))

	text := bitops.MultiXOR32(encB, chunks)

	if err := stream.WriteBinaryStream(filePath, text); err != nil {
		golog.Errorf("Error prepending to binary stream: %v\n", err)
	}

}

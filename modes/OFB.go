package modes

import (
	"lea/bitops"
	"lea/core"
	"lea/stream"

	"github.com/kopytkg/golog"
)

func encryptOFB(filePath string, prev *[4]uint32, keySegments []uint32, chunks [4]uint32, size int) {
	*prev = [4]uint32(core.SelectEncrypt(*prev, keySegments, size))

	encB := bitops.MultiXOR32(*prev, chunks)

	if err := stream.WriteBinaryStream(filePath, encB); err != nil {
		golog.Errorf("Error writing to binary stream: %v\n", err)
	}
}

func decryptOFB(filePath string, prev *[4]uint32, keySegments []uint32, chunks [4]uint32, size int) {
	encryptOFB(filePath, prev, keySegments, chunks, size)
}

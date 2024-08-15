package modes

import (
	"bufio"
	"fmt"
	"lea/bitops"
	"lea/core"
	"lea/errors"
	"lea/fingerprint"
	"lea/schedule"
	"lea/stream"
	"os"
)

func PerformOFB(filePath string, bKey []byte, bSeed []byte, encrypt bool, keySize int) {
	chunks := stream.BinaryChunkStream(filePath)
	
	kChunks := fingerprint.LoadSource(bKey)
	sChunks := fingerprint.LoadSource(bSeed)
	
        key := fingerprint.SelectPrint(kChunks, keySize)
        seed := fingerprint.SelectPrint(sChunks, keySize)

	rk := schedule.KeySchedule(keySize, key, seed)

	var blocks [4]uint32
	
	if encrypt {
		encryptOFB(filePath, blocks, rk, chunks, keySize)
	} else {
		decryptOFB(filePath, blocks, rk, chunks, keySize)
	}
	"fmt"
	"lea/bitops"
	"lea/core"
	"lea/stream"
)


func encryptOFB(filePath string, prev *[4]uint32, keySegments []uint32, chunks [4]uint32, size int) {
	*prev = [4]uint32(core.SelectEncrypt(*prev, keySegments, size))
	encB := bitops.MultiXOR32(*prev, chunks)
	if err := stream.WriteBinaryStream(filePath, encB); err != nil {
		fmt.Printf("Error writing to binary stream: %v\n", err)
	}
}

func encryptOFB(filePath string, blocks [4]uint32, keySegments []uint32, chunks []uint32, size int) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Encrypting", filePath)	
	var encChunks []uint32
	// ask user for IV through stdin
	fmt.Println("Please provide an IV: ")
	input, _ := reader.ReadString('\n')
	input = input[:len(input)-1]
	
	// conver the input to a slice of uint32
	var IV [4]uint32 = fingerprint.Fingerprint128(fingerprint.LoadSource([]byte(input)))

	var prev [4]uint32 = IV
	for i := 0; i < len(chunks); i++ {
		blocks[i%4] = chunks[i]
		if (i+1)%4 == 0 {
			prev = [4]uint32(core.SelectEncrypt(prev, keySegments, size))
			encryptedBlock := bitops.MultiXOR32(prev, blocks)
			encChunks = append(encChunks, encryptedBlock[:]...)
		}
	}

	if len(chunks)%4 != 0 {
		errors.PaddingError()
	}
	stream.WriteBinaryStream(filePath, encChunks)
}

func decryptOFB(filePath string, blocks [4]uint32, keySegments []uint32, chunks []uint32, size int) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Encrypting", filePath)	
	var encChunks []uint32
	// ask user for IV through stdin
	fmt.Println("Please provide an IV: ")
	input, _ := reader.ReadString('\n')
	input = input[:len(input)-1]
	
	

	// conver the input to a slice of uint32
	var IV [4]uint32 = fingerprint.Fingerprint128(fingerprint.LoadSource([]byte(input)))

	var prev [4]uint32 = IV	
	for i := 0; i < len(chunks); i++ {
		blocks[i%4] = chunks[i]
		if (i+1)%4 == 0 {
			prev = [4]uint32(core.SelectEncrypt(prev, keySegments, size))
			encryptedBlock := bitops.MultiXOR32(prev, blocks)
			encChunks = append(encChunks, encryptedBlock[:]...)
		}
	}

	if len(chunks)%4 != 0 {
		errors.PaddingError()
	}
	stream.WriteBinaryStream(filePath, encChunks)
}
func decryptOFB(filePath string, prev *[4]uint32, keySegments []uint32, chunks [4]uint32, size int) {
	encryptOFB(filePath, prev, keySegments, chunks, size)
}


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
	"lea/utils"
	"os"
)

func PerformCBC(filePath string, bKey []byte, bSeed []byte, encrypt bool, keySize int) {
	chunks := stream.BinaryChunkStream(filePath)
	
	kChunks := fingerprint.LoadSource(bKey)
	sChunks := fingerprint.LoadSource(bSeed)
	
        key := fingerprint.SelectPrint(kChunks, keySize)
        seed := fingerprint.SelectPrint(sChunks, keySize)

	rk := schedule.KeySchedule(keySize, key, seed)

	var blocks [4]uint32
	
	if encrypt {
		encryptCBC(filePath, blocks, rk, chunks, keySize)
	} else {
		decryptCBC(filePath, blocks, rk, chunks, keySize)
	}

}

func encryptCBC(filePath string, blocks [4]uint32, keySegments []uint32, chunks []uint32, size int) {
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
			after :=  bitops.MultiXOR32(blocks, prev)
			encryptedBlock := core.SelectEncrypt(after, keySegments, size)
			encChunks = append(encChunks, encryptedBlock[:]...)
			prev = [4]uint32(encryptedBlock)
		}
	}
	
	if len(chunks)%4 != 0 {
		errors.PaddingError()
	}
	stream.WriteBinaryStream(filePath, encChunks)
}



func decryptCBC(filePath string, blocks [4]uint32, keySegments []uint32, chunks []uint32, size int) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Encrypting", filePath)	
	var encChunks []uint32
	// ask user for IV through stdin
	fmt.Println("Please provide an IV: ")
	input, _ := reader.ReadString('\n')
	input = input[:len(input)-1]
	
	

	// conver the input to a slice of uint32
	var IV [4]uint32 = fingerprint.Fingerprint128(fingerprint.LoadSource([]byte(input)))
	
	var stackedBlocks utils.Stack = utils.Stack{}
	var flushStack utils.Stack = utils.Stack{}

	for i := 0; i < len(chunks); i++ {
		blocks[i%4] = chunks[i]
		if (i+1)%4 == 0 {
			stackedBlocks.Append(blocks)		
		}
	}
	for stackedBlocks.Length() > 0 {
		currentBlock := stackedBlocks.Pop()
		decryptedBlock := core.SelectDecrypt(currentBlock, keySegments, size)
		next := IV
		if stackedBlocks.Length() != 0 {next= stackedBlocks.Peek()}
		base := bitops.MultiXOR32([4]uint32(decryptedBlock), next)
		flushStack.Append(base)
	}
	for flushStack.Length() > 0 {
		el := flushStack.Pop()
		encChunks = append(encChunks, el[:]...)
	}

	if len(chunks)%4 != 0 {
		errors.PaddingError()
	}
	stream.WriteBinaryStream(filePath, encChunks)
}


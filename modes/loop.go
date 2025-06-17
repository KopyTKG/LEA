package modes

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"lea/fingerprint"
	"lea/schedule"
	"lea/state"
	"lea/stream"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
)

func PerformMode(encrypt bool) {

	var IV [4]uint32

	if state.CYPHERMODE != "ecb" {
		cli := bufio.NewReader(os.Stdin)
		fmt.Print("Please provide an IV: ")
		input, _ := cli.ReadString('\n')
		input = strings.TrimSpace(input)

		IV = fingerprint.Fingerprint128(fingerprint.LoadSource([]byte(input)))
	}

	var files []string
	var tmpFiles []string
	if state.RECURSION {
		content, err := stream.RecursionLS(state.FILEPATH)
		if err != nil {
			panic(err)
		}

		for _, file := range content {
			files = append(files, file)
			tmpFiles = append(tmpFiles, file+".tmp")
		}
	} else {
		files = append(files, state.FILEPATH)
		tmpFiles = append(tmpFiles, state.FILEPATH+".tmp")
	}

	if state.VERBOSE {

	}

	fmt.Println(files)

	maxWorkers := runtime.NumCPU()

	semaphore := make(chan struct{}, maxWorkers)

	var wg sync.WaitGroup

	wg.Add(len(files))

	for i := 0; i < len(files); i++ {
		fmt.Println(i)
		t := IV
		go worker(&wg, semaphore, tmpFiles[i], files[i], encrypt, &t)
	}

}

func worker(wg *sync.WaitGroup, semaphore chan struct{}, tmp, path string, enc bool, IV *[4]uint32) {

	fmt.Println("Worker on file " + path)
	defer wg.Done()

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Skipping " + path)
		return
	}
	defer file.Close()

	kChunks := fingerprint.LoadSource(state.ByteKEY)
	sChunks := fingerprint.LoadSource(state.ByteSEED)
	key := fingerprint.SelectPrint(kChunks, state.KEYLENGTH)
	seed := fingerprint.SelectPrint(sChunks, state.KEYLENGTH)
	rk := schedule.KeySchedule(state.KEYLENGTH, key, seed)

	semaphore <- struct{}{}

	defer func() { <-semaphore }()

	fmt.Println(path)
	readAndProcessFileInChunks(state.CYPHERMODE, tmp, rk, file, IV, enc, state.KEYLENGTH)

	cleanup(tmp, path)

}

func readAndProcessFileInChunks(mode string, tmpFilePath string, rk []uint32, file *os.File, prev *[4]uint32, encrypt bool, keySize int) {
	reader := bufio.NewReader(file)
	var chunks []uint32
	buf := make([]byte, state.CHUNKSIZE)
	count := 0
	for {
		n, err := io.ReadFull(reader, buf)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			if n > 0 {
				var paddedBuf [state.CHUNKSIZE]byte
				copy(paddedBuf[:], buf[:n])
				chunks = append(chunks, binary.LittleEndian.Uint32(paddedBuf[:]))
			}
			break
		}
		if err != nil {
			log.Fatalf("Error reading file: %v", err)
		}

		chunk := binary.LittleEndian.Uint32(buf)
		chunks = append(chunks, chunk)

		if len(chunks) == 4 {
			performAction(mode, tmpFilePath, rk, [4]uint32(chunks), prev, encrypt, keySize)
			chunks = []uint32{}
			count += 16

		}

	}

	for len(chunks)%4 != 0 {
		chunks = append(chunks, 0)
	}
	if len(chunks) > 0 {
		performAction(mode, tmpFilePath, rk, [4]uint32(chunks), prev, encrypt, keySize)
	}
}

func cleanup(tmp, filePath string) {
	if err := os.Remove(filePath); err != nil {
		log.Fatalf("Error removing original file: %v", err)
	}
	if err := os.Rename(tmp, filePath); err != nil {
		log.Fatalf("Error renaming temporary file: %v", err)
	}
}

func performAction(mode, filePath string, rk []uint32, chunks [4]uint32, prev *[4]uint32, encrypt bool, keySize int) {
	switch mode {
	default:
		log.Fatalln("No mode selected")
		os.Exit(1)

	case "ecb":
		if encrypt {
			encryptECB(filePath, rk, chunks, keySize)
		} else {
			decryptECB(filePath, rk, chunks, keySize)
		}

	case "cbc":
		if encrypt {
			encryptCBC(filePath, prev, rk, chunks, keySize)
		} else {
			decryptCBC(filePath, prev, rk, chunks, keySize)
		}

	case "cfb":
		if encrypt {
			encryptCFB(filePath, prev, rk, chunks, keySize)
		} else {
			decryptCFB(filePath, prev, rk, chunks, keySize)
		}

	case "ofb":
		if encrypt {
			encryptOFB(filePath, prev, rk, chunks, keySize)
		} else {
			decryptOFB(filePath, prev, rk, chunks, keySize)
		}
	}
}

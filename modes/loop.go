package modes

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"lea/fingerprint"
	"lea/schedule"
	"lea/terminal"
	"log"
	"os"
	"strings"
	"time"

	ui "github.com/gizak/termui/v3"
)

const chunkSize = 4

func PerformMode(mode, filePath string, bKey, bSeed []byte, encrypt bool, keySize int) {
	kChunks := fingerprint.LoadSource(bKey)
	sChunks := fingerprint.LoadSource(bSeed)
	key := fingerprint.SelectPrint(kChunks, keySize)
	seed := fingerprint.SelectPrint(sChunks, keySize)
	rk := schedule.KeySchedule(keySize, key, seed)
	tmpFilePath := filePath + ".tmp"

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	var prev [4]uint32
	var IV [4]uint32

	if mode != "ecb" {
		cli := bufio.NewReader(os.Stdin)
		fmt.Print("Please provide an IV: ")
		input, _ := cli.ReadString('\n')
		input = strings.TrimSpace(input)

		IV = fingerprint.Fingerprint128(fingerprint.LoadSource([]byte(input)))
	}

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	defer ui.Close()

	fs, err := file.Stat()
	size := int(fs.Size())
	w, _ := ui.TerminalDimensions()
	bar := terminal.BarSetup((w - 5) / 2)
	f := terminal.Fileln{
		FP:    filePath,
		Done:  0,
		Total: size,
		Bar:   bar,
	}
	r := terminal.Rendering{File: &f}

	prev = IV
	readAndProcessFileInChunks(mode, tmpFilePath, rk, file, &prev, encrypt, keySize, &f, &r)

	cleanup(tmpFilePath, filePath)
}

func readAndProcessFileInChunks(mode string, tmpFilePath string, rk []uint32, file *os.File, prev *[4]uint32, encrypt bool, keySize int, f *terminal.Fileln, r *terminal.Rendering) {
	reader := bufio.NewReader(file)
	var chunks []uint32
	buf := make([]byte, chunkSize)
	count := 0
	lastRenderTime := time.Now()
	renderInterval := 500 * time.Millisecond // Set the interval for UI updates
	for {
		n, err := io.ReadFull(reader, buf)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			if n > 0 {
				var paddedBuf [chunkSize]byte
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
			f.Update(count)
		}
		if time.Since(lastRenderTime) > renderInterval {
			r.Run()
			lastRenderTime = time.Now()
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

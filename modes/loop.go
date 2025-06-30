package modes

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"lea/core"
	"lea/schedule"
	"lea/state"
	"lea/stream"
	"lea/terminal"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/kopytkg/golog"
)

func PerformMode() {
	rk, err := schedule.KeySchedule(int(state.Key.Metadata.KeyLength), state.Key.Key, state.Key.Seed)

	if err != nil {
		golog.Errorf("Error in key schedule: %v", err)
		os.Exit(1)
	}

	var files []string

	UI := terminal.Rendering{}
	UI.SetupOS()
	UI.Files = make([]*terminal.Fileln, 0)

	if state.RECURSION {
		content, err := stream.RecursionLS(state.FILEPATH)
		if err != nil {
			panic(err)
		}

		for _, file := range content {
			files = append(files, file)
		}
	} else {
		files = append(files, state.FILEPATH)
	}

	if state.VERBOSE {

	}

	maxWorkers := runtime.NumCPU()

	semaphore := make(chan struct{}, maxWorkers)

	UI.Total = len(files)

	var wg sync.WaitGroup

	wg.Add(len(files))

	for i := 0; i < len(files); i++ {
		go worker(&wg, semaphore, files[i], rk, &UI)
	}

	if state.VERBOSE {
		go func() {
			for {
				UI.Run()
				time.Sleep(200 * time.Millisecond)
			}
		}()
	}

	wg.Wait()

}

func worker(wg *sync.WaitGroup, semaphore chan struct{}, path string, rk []uint32, UI *terminal.Rendering) {
	defer wg.Done()

	// Acquire semaphore slot immediately to respect concurrency limit
	semaphore <- struct{}{}
	defer func() { <-semaphore }()

	// Create Target
	target := core.Target{}

	if err := target.Open(path); err != nil {
		golog.Errorf("Error opening file (%s): %v", path, err)
		return
	}
	fs, _ := target.File.Stat()

	f := terminal.Fileln{Filename: path, Total: uint64(fs.Size()), Current: 0, Done: false}

	(*UI).AddFile(&f)
	f.Bar = terminal.BarSetup(50)

	metadata := core.Metadata{}

	if !state.ENCRYPT {
		mete, err := core.ReadMetadata(&target)
		if err != nil {
			golog.Error(err)
			return
		}

		metadata = mete
	} else {
		if state.CYPHERMODE != "ecb" {
			if err := target.RandomIV(); err != nil {
				golog.Errorf("Error generating random IV: %v", err)
				return
			}
		} else {
			for i := range target.IV {
				target.IV[i] = 0
			}

		}
	}

	var args ModeArgs

	var marker byte = '-'

	switch state.CYPHERMODE {
	case "ecb":
		args = &ECBArgs{FilePath: target.Temp, RK: rk, KeySize: int(state.Key.Metadata.KeyLength)}
		marker = 'E'
	case "cbc":
		args = &CBCArgs{FilePath: target.Temp, RK: rk, KeySize: int(state.Key.Metadata.KeyLength), Prev: &target.IV}
		marker = 'B'
	case "cfb":
		args = &CFBArgs{FilePath: target.Temp, RK: rk, KeySize: int(state.Key.Metadata.KeyLength), Prev: &target.IV}
		marker = 'F'
	case "ofb":
		args = &OFBArgs{FilePath: target.Temp, RK: rk, KeySize: int(state.Key.Metadata.KeyLength), Prev: &target.IV}
		marker = 'O'
	case "ctr":
		args = &CTRArgs{FilePath: target.Temp, RK: rk, KeySize: int(state.Key.Metadata.KeyLength), Counter: &target.IV}
		marker = 'T'
	default:
		golog.Errorf("Unsupported cipher mode: %s", state.CYPHERMODE)
		return
	}

	if state.ENCRYPT {
		metadata.IV = target.IV
		metadata.OriginalSize = fs.Size()
		metadata.Version = 1
		metadata.Mode = marker
		metadata.KeySize = uint16(state.Key.Metadata.KeyLength)
	}

	if state.ENCRYPT {
		bytes := make([]byte, 32)
		for i, v := range metadata.ToArray() {
			binary.LittleEndian.PutUint32(bytes[i*4:(i+1)*4], v)
		}
		if err := stream.WriteBinaryStream(target.Temp, bytes); err != nil {
			golog.Errorf("Error writing to binary stream: %v\n", err)
			return
		}
	}

	var bytesWritten uint64 = 0
	for chunk := range target.Stream() {
		block, err := CryptMethod(chunk, args)
		if err != nil {
			golog.Errorf("Error during encryption/decryption: %v", err)
			return
		}
		f.Update(bytesWritten)

		bytes := make([]byte, 16)
		for i, v := range block {
			binary.LittleEndian.PutUint32(bytes[i*4:(i+1)*4], v)
		}

		toWrite := 16
		if !state.ENCRYPT {
			// Only on the last block, write up to original size
			remaining := int64(metadata.OriginalSize) - int64(bytesWritten)
			if remaining < 16 {
				if remaining <= 0 {
					break
				}
				toWrite = int(remaining)
			}
		}

		if err := stream.WriteBinaryStream(target.Temp, bytes[:toWrite]); err != nil {
			golog.Errorf("Error writing to binary stream: %v\n", err)
			return
		}

		bytesWritten += uint64(toWrite)

		if !state.ENCRYPT && bytesWritten >= uint64(metadata.OriginalSize) {
			break
		}
	}

	target.File.Close()
	cleanup(target.Temp, target.Host)

	f.Done = true
	UI.Done += 1
}

func makeRandomBuffer(size int64) ([]byte, error) {
	buf := make([]byte, size)

	if _, err := rand.Read(buf); err != nil {
		return nil, err
	}

	return buf, nil
}

func secureWipe(path string, iterations int) error {
	file, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return err
	}

	stat, _ := file.Stat()
	size := stat.Size()

	buf, err := makeRandomBuffer(size)
	if err != nil {
		golog.Error(err)
		os.Exit(1)
	}

	patterns := [][]byte{
		bytes.Repeat([]byte{0x00}, int(size)),
		bytes.Repeat([]byte{0xFF}, int(size)),
		buf,
	}

	for i := range iterations {
		pattern := patterns[i%len(patterns)]
		if _, err := file.WriteAt(pattern, 0); err != nil {
			return err
		}
		file.Sync()
	}

	file.Close()
	return os.Remove(path)
}

func cleanup(tmp, filePath string) {

	if err := secureWipe(filePath, state.Iterations); err != nil {
		golog.Errorf("Error removing original file: %v", err)
	}

	if err := os.Rename(tmp, filePath); err != nil {
		golog.Errorf("Error renaming temporary file: %v", err)
	}
}

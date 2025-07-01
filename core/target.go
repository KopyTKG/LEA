package core

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"iter"
	"lea/state"
	"os"

	"github.com/kopytkg/golog"
)

// Metadata holds metadata information for decryption operation
type Metadata struct {
	IV           [4]uint32 // Initialization Vector used for encryption (counter for CTR, IV for CBC, etc.)
	OriginalSize int64     // Original size of the file before encryption
	Version      uint8     // Version of the encryption algorithm used
	Mode         byte      // 'O' for OFB, 'T' for CTR, 'E' for ECB, 'F' for CFB, 'B' for CBC
	KeySize      uint16    // Size of the key used for encryption (128, 192, or 256 bits)
	CMAC         [4]uint32 // Optional: HMAC or other integrity check data
}

func (m *Metadata) ToArray() []uint32 {
	// Convert Metadata to a []uint32 array for writing
	return []uint32{
		m.IV[0], m.IV[1], m.IV[2], m.IV[3],
		uint32(m.OriginalSize),
		uint32(m.Version),
		uint32(m.Mode),
		uint32(m.KeySize),
		m.CMAC[0], m.CMAC[1], m.CMAC[2], m.CMAC[3],
	}
}

const (
	OFFSET = 48 // Offset for metadata in the file (16 bytes for IV + 16 bytes for metadata + 16 bytes for CMAC)
)

// Target represents a file target for encryption/decryption operations
type Target struct {
	Host string
	Temp string

	File *os.File

	IV [4]uint32

	OnChunk func(chunk [4]uint32)
}

func ReadMetadata(t *Target) (Metadata, error) {
	if t.File == nil {
		return Metadata{}, fmt.Errorf("File is not opened, cannot read metadata")

	}

	meta := Metadata{}

	buf := make([]byte, 48) // 4 uint32 values = 16 bytes IV + 16 bytes for metadata + 16 bytes for CMAC

	_, err := io.ReadFull(t.File, buf)
	if err != nil {
		return Metadata{}, fmt.Errorf("failed to read metadata: %w", err)
	}

	// Separate the first 16 bytes for IV
	for i := range 4 {
		t.IV[i] = binary.LittleEndian.Uint32(buf[i*4 : (i+1)*4])
	}

	meta.IV = t.IV

	// The rest of the buffer can be used for metadata
	metadata := buf[16:32]
	if len(metadata) < 16 {
		return Metadata{}, fmt.Errorf("metadata buffer is too small")
	}

	for i := range 4 {
		switch i {
		case 0:
			meta.OriginalSize = int64(binary.LittleEndian.Uint32(metadata[i*4 : (i+1)*4])) // Original size
		case 1:
			meta.Version = uint8(binary.LittleEndian.Uint32(metadata[i*4 : (i+1)*4])) // Version
		case 2:
			meta.Mode = byte(binary.LittleEndian.Uint32(metadata[i*4 : (i+1)*4])) // Mode
		case 3:
			meta.KeySize = uint16(binary.LittleEndian.Uint32(metadata[i*4 : (i+1)*4])) // Key size
		}
	}

	// The last 16 bytes can be used for CMAC or other integrity checks
	meta.CMAC[0] = binary.LittleEndian.Uint32(buf[32:36])
	meta.CMAC[1] = binary.LittleEndian.Uint32(buf[36:40])
	meta.CMAC[2] = binary.LittleEndian.Uint32(buf[40:44])
	meta.CMAC[3] = binary.LittleEndian.Uint32(buf[44:48])

	return meta, nil
}

func (t *Target) RandomIV() error {
	buf := make([]byte, 16)

	if _, err := rand.Read(buf); err != nil {
		return err
	}

	for i := range 4 {
		t.IV[i] = binary.LittleEndian.Uint32(buf[i*4 : (i+1)*4])
	}

	return nil
}

func (t *Target) Open(path string) error {
	file, err := os.Open(path)

	if err != nil {
		return err
	}

	t.Host = path
	t.Temp = path + ".tmp"

	t.File = file
	return nil
}

func (t *Target) Close() error {
	if t.File != nil {
		err := t.File.Close()
		if err != nil {
			return err
		}
		t.File = nil
	}

	t.Host = ""
	t.Temp = ""

	return nil
}

func (t *Target) Stream() iter.Seq[[4]uint32] {
	return func(yield func([4]uint32) bool) {
		var chunks []uint32
		buf := make([]byte, state.CHUNKSIZE)

		for {
			n, err := io.ReadFull(t.File, buf)
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				if n > 0 {
					var paddedBuf [state.CHUNKSIZE]byte
					copy(paddedBuf[:], buf[:n])
					chunks = append(chunks, binary.LittleEndian.Uint32(paddedBuf[:]))
				}
				break
			}
			if err != nil {
				golog.Errorf("Error reading file: %v", err)
				return
			}

			chunk := binary.LittleEndian.Uint32(buf)
			chunks = append(chunks, chunk)

			if len(chunks) == 4 {
				var out [4]uint32
				copy(out[:], chunks)
				if !yield(out) {
					return
				}
				chunks = chunks[:0]
			}
		}

		if len(chunks) > 0 {
			for len(chunks) < 4 {
				chunks = append(chunks, 0)
			}
			var out [4]uint32
			copy(out[:], chunks)
			if !yield(out) {
				return
			}
		}
	}
}

func (t *Target) ResetStream() {
	// Reset the file pointer to the header OFFSET
	if t.File != nil {
		if _, err := t.File.Seek(OFFSET, io.SeekStart); err != nil {
			golog.Errorf("Error resetting file stream: %v", err)
		}
	} else {
		golog.Error("File is not opened, cannot reset stream")
	}
}

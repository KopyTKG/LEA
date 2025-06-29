package core

import (
	"bufio"
	"encoding/binary"
	"io"
	"iter"
	"lea/state"
	"os"

	"github.com/kopytkg/golog"
)

type Target struct {
	Host string
	Temp string

	File *os.File

	IV [4]uint32

	OnChunk func(chunk [4]uint32)
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
		reader := bufio.NewReader((*t).File)
		var chunks []uint32
		buf := make([]byte, state.CHUNKSIZE)

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

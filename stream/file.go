package stream


import (
import (
	"bufio"
	"io"
	"log"
	"os"
)

func GetFile(path string) []byte {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("File (%v) could not be accessed\n\n", path)
		os.Exit(1)
		return []byte{}
	} else {
		return BinaryStream(path)
		file, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		reader := bufio.NewReader(file)
		var chunks []byte
		for {
			n, err := reader.ReadByte()
			if err != nil {
				if err == io.EOF {
					break
				}
			}
			chunks = append(chunks, n)
		}
		return chunks
	}
}

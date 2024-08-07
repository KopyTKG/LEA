package stream


import (
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
	}
}

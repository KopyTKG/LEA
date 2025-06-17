package stream

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

func LSFolder(path string) ([]string, error) {
	entries, err := os.ReadDir(path)

	if err != nil {
		return []string{}, err
	}

	var paths []string

	for _, entry := range entries {
		p := path + "/" + entry.Name()
		paths = append(paths, p)
	}

	return paths, nil
}

func IsFolder(path string) (bool, error) {
	file, err := os.Open(path)

	if err != nil {
		return false, err
	}

	defer file.Close()

	fileInfo, err := file.Stat()

	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), nil
}

func RecursionLS(path string) ([]string, error) {
	c, err := LSFolder(path)
	if err != nil {
		return []string{}, err
	}
	var files []string
	for _, f := range c {
		b, err := IsFolder(f)

		if err != nil {
			continue
		}

		if !b {
			files = append(files, f)
		} else {
			fs, err := RecursionLS(f)
			if err != nil {
				return []string{}, err
			}
			files = append(files, fs...)
		}
	}

	return files, nil
}

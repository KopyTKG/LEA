package stream

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func GetFile(path string) ([]byte, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return []byte{}, fmt.Errorf("File (%v) could not be accessed", path)
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
		return chunks, nil
	}
}

func LSFolder(path string) ([]string, error) {
	entries, err := os.ReadDir(path)

	if err != nil {
		return []string{}, fmt.Errorf("Folder (%v) could not be accessed: %w", path, err)
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
		return false, fmt.Errorf("File (%v) could not be accessed: %w", path, err)
	}

	defer file.Close()

	fileInfo, err := file.Stat()

	if err != nil {
		return false, fmt.Errorf("Could not get file info for (%v): %w", path, err)
	}

	return fileInfo.IsDir(), nil
}

func RecursionLS(path string) ([]string, error) {
	c, err := LSFolder(path)
	if err != nil {
		return []string{}, fmt.Errorf("Could not list folder (%v): %w", path, err)
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

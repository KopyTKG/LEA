package key

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kopytkg/golog"
)

func ReadArmoredKey(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open key file: %w", err)
	}

	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read key file: %w", err)
	}
	// Convert content to string and trim whitespace
	armored := strings.TrimSpace(string(content))
	return armored, nil
}

func ParseArmoredKey(armored string) (*KeyMaterial, error) {
	lines := strings.Split(armored, "\n")
	km := &KeyMaterial{}
	var keyHex, seedHex string

	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "Version: "):
			v := strings.TrimPrefix(line, "Version: ")
			vint, err := strconv.Atoi(v)
			if err != nil {
				golog.Error(err)
				panic("")
			}
			km.Version = uint(vint)

		case strings.HasPrefix(line, "Created: "):
			km.Metadata.Created = strings.TrimPrefix(line, "Created: ")
		case strings.HasPrefix(line, "Key Length: "):
			fmt.Sscanf(line, "Key Length: %d bits", &km.Metadata.KeyLength)
		case strings.HasPrefix(line, "Key: "):
			keyHex = strings.ReplaceAll(strings.TrimPrefix(line, "Key: "), ":", "")
		case strings.HasPrefix(line, "Key Digest: "):
			km.Metadata.KeyDigest = strings.TrimPrefix(line, "Key Digest: ")
		case strings.HasPrefix(line, "Seed: "):
			seedHex = strings.ReplaceAll(strings.TrimPrefix(line, "Seed: "), ":", "")
		case strings.HasPrefix(line, "Seed Digest: "):
			km.Metadata.SeedDigest = strings.TrimPrefix(line, "Seed Digest: ")
		}
	}

	// Convert hex to uint32 arrays
	if err := hexToUint32(keyHex, &km.Key); err != nil {
		return nil, fmt.Errorf("invalid key format: %w", err)
	}
	if seedHex != "" {
		if err := hexToUint32(seedHex, &km.Seed); err != nil {
			return nil, fmt.Errorf("invalid seed format: %w", err)
		}
	}

	// Validate digests
	if err := km.Validate(); err != nil {
		return nil, err
	}

	// Validate timestamp
	if _, err := time.Parse(time.RFC3339, km.Metadata.Created); err != nil {
		return nil, fmt.Errorf("invalid creation time: %w", err)
	}

	return km, nil
}

func hexToUint32(hexStr string, dest *[]uint32) error {
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return err
	}

	if len(data)%4 != 0 {
		return fmt.Errorf("hex length must be multiple of 8 (got %d)", len(hexStr))
	}

	*dest = make([]uint32, len(data)/4)
	for i := range *dest {
		(*dest)[i] = binary.BigEndian.Uint32(data[i*4:])
	}
	return nil
}

package key

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"lea/types"
	"time"
)

func KeyGen(size int) (*KeyMaterial, error) {
	key, err := getKey(size)
	if err != nil {
		return nil, err
	}

	seed, err2 := getSeed(size)
	if err2 != nil {
		return nil, err2
	}

	var km KeyMaterial = KeyMaterial{}

	km.Version = 1
	km.Key = key
	km.Seed = seed
	km.Metadata.Created = time.Now().Format(time.RFC3339)
	km.Metadata.KeyLength = uint16(size)

	// Calculate key digest
	dkey, dseed := km.Digests()
	km.Metadata.KeyDigest = dkey
	km.Metadata.SeedDigest = dseed
	return &km, nil
}

func getKey(size int) ([]uint32, error) {
	switch size {
	case 128:
		chunks, err := chunks128()
		if err != nil {
			return nil, fmt.Errorf("error generating 128-bit key: %w", err)
		}
		return chunks[:], nil
	case 192:
		chunks, err := chunks192()
		if err != nil {
			return nil, fmt.Errorf("error generating 192-bit key: %w", err)
		}
		return chunks[:], nil
	case 256:
		chunks, err := chunks256()
		if err != nil {
			return nil, fmt.Errorf("error generating 256-bit key: %w", err)
		}
		return chunks[:], nil

	default:
		return []uint32{}, fmt.Errorf("invalid key size %d, must be 128, 192 or 256", size)
	}
}

func getSeed(size int) ([]uint32, error) {
	switch size {
	case 128:
		chunks, err := chunks128()
		if err != nil {
			return nil, fmt.Errorf("error generating 128-bit seed: %w", err)
		}
		return chunks[:], nil
	case 192:
		chunks, err := chunks192()
		if err != nil {
			return nil, fmt.Errorf("error generating 192-bit seed: %w", err)
		}
		return chunks[:], nil
	case 256:
		chunks, err := chunks256()
		if err != nil {
			return nil, fmt.Errorf("error generating 256-bit seed: %w", err)
		}
		return chunks[:], nil

	default:
		return []uint32{}, fmt.Errorf("invalid seed size %d, must be 128, 192 or 256", size)
	}
}

func chunks128() (*types.Key128, error) {
	var chunks types.Key128

	buf := make([]byte, 16)

	if _, err := rand.Read(buf); err != nil {
		return nil, err
	}

	for i := range 4 {
		chunks[i] = binary.BigEndian.Uint32(buf[i*4 : (i+1)*4])
	}

	return &chunks, nil
}

func chunks192() (*types.Key192, error) {
	var chunks types.Key192

	buf := make([]byte, 24)

	if _, err := rand.Read(buf); err != nil {
		return nil, err
	}

	for i := range 6 {
		chunks[i] = binary.BigEndian.Uint32(buf[i*4 : (i+1)*4])
	}

	return &chunks, nil
}

func chunks256() (*types.Key256, error) {
	var chunks types.Key256

	buf := make([]byte, 32)

	if _, err := rand.Read(buf); err != nil {
		return nil, err
	}

	for i := range 8 {
		chunks[i] = binary.BigEndian.Uint32(buf[i*4 : (i+1)*4])
	}

	return &chunks, nil
}

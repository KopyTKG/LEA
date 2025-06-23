package key

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/kopytkg/golog"
)

type KeyMaterial struct {
	Version  uint     `json:"version"` // Version of the key file format
	Key      []uint32 `json:"key"`     // Key in 32-bit chunks
	Seed     []uint32 `json:"seed"`    // Seed in 32-bit chunks
	Metadata struct {
		Created    string `json:"created"`    // ISO 8601 format
		KeyLength  uint16 `json:"keyLength"`  // Length in bits
		KeyDigest  string `json:"keyDigest"`  // SHA-256 hash of the key
		SeedDigest string `json:"seedDigest"` // SHA-256 hash of the seed
	} `json:"meta"`
}

func (km *KeyMaterial) Digests() (string, string) {
	key := computeDigest(km.Key)
	seed := computeDigest(km.Seed)
	return key, seed
}

// Computes SHA-256 digest for 128/192/256-bit keys (4/6/8 uint32 values)
func computeDigest(chunks []uint32) string {
	switch len(chunks) {
	case 4, 6, 8: // Validate expected lengths
		buf := make([]byte, len(chunks)*4)
		for i, v := range chunks {
			binary.BigEndian.PutUint32(buf[i*4:], v)
		}
		hash := sha256.Sum256(buf)
		return hex.EncodeToString(hash[:])
	default:
		golog.Error("invalid key length - expected 4/6/8 uint32 values")
		os.Exit(1)
		return ""
	}
}

func (km *KeyMaterial) Validate() error {
	// Check version
	if km.Version != 1 {
		return fmt.Errorf("unsupported version %d", km.Version)
	}

	// Validate key length
	expectedLenght := int(km.Metadata.KeyLength) / 8
	if len(km.Key)*4 != expectedLenght {
		return fmt.Errorf("key length mismatch: expected %d bytes, got %d", expectedLenght/8, len(km.Key)*4)
	}

	// Verify digests
	if err := km.VerifyDigest(); err != nil {
		return err
	}

	return nil
}

func (km *KeyMaterial) VerifyDigest() error {
	dkey, dseed := km.Digests()

	if dkey != km.Metadata.KeyDigest {
		return fmt.Errorf("key digest verification failed")
	}

	if dseed != km.Metadata.SeedDigest {
		return fmt.Errorf("seed digest verification failed")
	}

	return nil

}

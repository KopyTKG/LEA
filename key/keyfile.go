package key

import (
	"fmt"
	"os"
	"strings"

	"github.com/kopytkg/golog"
)

func KeyFileGen(km KeyMaterial, path string) error {
	// Validate key material
	if err := km.Validate(); err != nil {
		return fmt.Errorf("invalid key material: %w", err)
	}

	// Create file with restrictive permissions (0400)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0400)
	if err != nil {
		return fmt.Errorf("failed to create key file: %w", err)
	}
	defer file.Close()

	// Generate human-readable output
	output := generateKeyFileOutput(km)

	// Write to file
	if _, err := file.WriteString(output); err != nil {
		return fmt.Errorf("failed to write key file: %w", err)
	}

	golog.Infof("Successfully generated key file: %s", path)
	return nil
}

func generateKeyFileOutput(km KeyMaterial) string {
	var builder strings.Builder

	builder.WriteString("-----BEGIN LEA KEY-----\n")
	builder.WriteString(fmt.Sprintf("Version: %d\n", km.Version))
	builder.WriteString(fmt.Sprintf("Created: %s\n", km.Metadata.Created))
	builder.WriteString(fmt.Sprintf("Key Length: %d bits\n", km.Metadata.KeyLength))

	builder.WriteString("Key: ")
	for i, chunk := range km.Key {
		if i > 0 {
			builder.WriteString(":")
		}
		builder.WriteString(fmt.Sprintf("%08X", chunk))
	}
	builder.WriteString(fmt.Sprintf("\nKey Digest: %s\n", km.Metadata.KeyDigest))

	if len(km.Seed) > 0 {
		builder.WriteString("Seed: ")
		for i, chunk := range km.Seed {
			if i > 0 {
				builder.WriteString(":")
			}
			builder.WriteString(fmt.Sprintf("%08X", chunk))
		}
		builder.WriteString(fmt.Sprintf("\nSeed Digest: %s\n", km.Metadata.SeedDigest))
	}

	builder.WriteString("-----END LEA KEY-----\n")
	return builder.String()
}

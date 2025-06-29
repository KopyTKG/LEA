package modes

import (
	"fmt"
	"lea/bitops"
	"lea/core"
	"lea/state"
	"lea/stream"
)

type CBC struct{}

type CBCArgs struct {
	FilePath string
	RK       []uint32
	KeySize  int
	Prev     *[4]uint32
}

func (c *CBC) argsCheck(args *CBCArgs) error {
	if args.RK == nil || len(args.RK) == 0 {
		return fmt.Errorf("Key schedule is nil or empty, cannot perform ECB encryption")
	}

	if args.FilePath == "" {
		return fmt.Errorf("File path is empty, cannot write encrypted data")
	}

	if !state.ValidKeys[args.KeySize] {
		return fmt.Errorf("Invalid key size: %d. Valid sizes are 128, 192, or 256 bits", args.KeySize)
	}

	if args.Prev == nil {
		return fmt.Errorf("Previous block is nil, cannot perform CBC encryption")
	}

	return nil
}

func (c *CBC) Encrypt(chunks [4]uint32, args ModeArgs) error {
	cbcArgs, ok := args.(*CBCArgs)
	if !ok {
		return fmt.Errorf("invalid args for CBC")
	}

	if err := c.argsCheck(cbcArgs); err != nil {
		return err
	}

	xorredBlock := bitops.MultiXOR32(chunks, *(cbcArgs.Prev))

	encryptedBlock := core.SelectEncrypt(xorredBlock, cbcArgs.RK, cbcArgs.KeySize)

	cbcArgs.Prev = (*[4]uint32)(encryptedBlock)

	if err := stream.WriteBinaryStream(cbcArgs.FilePath, [4]uint32(encryptedBlock)); err != nil {
		return fmt.Errorf("Error writing to binary stream: %v\n", err)
	}

	return nil
}

func (c *CBC) Decrypt(chunks [4]uint32, args ModeArgs) error {
	cbcArgs, ok := args.(*CBCArgs)
	if !ok {
		return fmt.Errorf("invalid args for CBC")
	}

	if err := c.argsCheck(cbcArgs); err != nil {
		return err
	}
	decryptedBlock := core.SelectDecrypt(chunks, cbcArgs.RK, cbcArgs.KeySize)

	xorredBlock := bitops.MultiXOR32([4]uint32(decryptedBlock), *(cbcArgs.Prev))

	cbcArgs.Prev = &chunks

	if err := stream.WriteBinaryStream(cbcArgs.FilePath, xorredBlock); err != nil {
		return fmt.Errorf("Error writing to binary stream: %v\n", err)
	}

	return nil
}

package modes

import (
	"fmt"
	"lea/bitops"
	"lea/core"
	"lea/state"
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

func (c *CBC) Encrypt(chunks [4]uint32, args ModeArgs) ([4]uint32, error) {
	cbcArgs, ok := args.(*CBCArgs)
	if !ok {
		return [4]uint32{}, fmt.Errorf("invalid args for CBC")
	}

	if err := c.argsCheck(cbcArgs); err != nil {
		return [4]uint32{}, err
	}

	xorredBlock := bitops.MultiXOR32(chunks, *(cbcArgs.Prev))

	encryptedBlock := core.SelectEncrypt(xorredBlock, cbcArgs.RK, cbcArgs.KeySize)

	*cbcArgs.Prev = encryptedBlock

	return encryptedBlock, nil
}

func (c *CBC) Decrypt(chunks [4]uint32, args ModeArgs) ([4]uint32, error) {
	cbcArgs, ok := args.(*CBCArgs)
	if !ok {
		return [4]uint32{}, fmt.Errorf("invalid args for CBC")
	}

	if err := c.argsCheck(cbcArgs); err != nil {
		return [4]uint32{}, err
	}

	decryptedBlock := core.SelectDecrypt(chunks, cbcArgs.RK, cbcArgs.KeySize)

	xorredBlock := bitops.MultiXOR32(decryptedBlock, *(cbcArgs.Prev))

	*cbcArgs.Prev = chunks

	return xorredBlock, nil
}

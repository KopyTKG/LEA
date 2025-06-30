package modes

import (
	"fmt"
	"lea/bitops"
	"lea/core"
	"lea/state"
)

type CFB struct{}

type CFBArgs struct {
	FilePath string
	RK       []uint32
	KeySize  int
	Prev     *[4]uint32
}

func (c *CFB) argsCheck(args *CFBArgs) error {
	if args.RK == nil || len(args.RK) == 0 {
		return fmt.Errorf("Key schedule is nil or empty, cannot perform CFB encryption")
	}

	if args.FilePath == "" {
		return fmt.Errorf("File path is empty, cannot write encrypted data")
	}

	if !state.ValidKeys[args.KeySize] {
		return fmt.Errorf("Invalid key size: %d. Valid sizes are 128, 192, or 256 bits", args.KeySize)
	}

	if args.Prev == nil {
		return fmt.Errorf("Previous block is nil, cannot perform CFB encryption")
	}

	return nil
}

func (c *CFB) Encrypt(chunks [4]uint32, args ModeArgs) ([4]uint32, error) {
	cfbArgs, ok := args.(*CFBArgs)
	if !ok {
		return [4]uint32{}, fmt.Errorf("invalid args for CFB")
	}

	if err := c.argsCheck(cfbArgs); err != nil {
		return [4]uint32{}, err
	}

	encryptedBlock := core.SelectEncrypt(*cfbArgs.Prev, cfbArgs.RK, cfbArgs.KeySize)

	xorredBlock := bitops.MultiXOR32(encryptedBlock, chunks)
	*cfbArgs.Prev = xorredBlock

	return xorredBlock, nil
}

func (c *CFB) Decrypt(chunks [4]uint32, args ModeArgs) ([4]uint32, error) {
	cfbArgs, ok := args.(*CFBArgs)
	if !ok {
		return [4]uint32{}, fmt.Errorf("invalid args for CFB")
	}

	if err := c.argsCheck(cfbArgs); err != nil {
		return [4]uint32{}, err
	}

	encryptedBlock := core.SelectEncrypt(*cfbArgs.Prev, cfbArgs.RK, cfbArgs.KeySize)

	xorredBlock := bitops.MultiXOR32(encryptedBlock, chunks)
	*cfbArgs.Prev = chunks

	return xorredBlock, nil
}

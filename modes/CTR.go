package modes

import (
	"fmt"
	"lea/bitops"
	"lea/core"
	"lea/state"
	"lea/stream"
)

type CTR struct{}

type CTRArgs struct {
	FilePath string
	RK       []uint32
	KeySize  int
	Counter  *[4]uint32
}

func (a *CTRArgs) incrementCounter() {
	for i := 3; i >= 0; i-- {
		a.Counter[i]++
		if a.Counter[i] != 0 {
			break
		}
	}
}

func (c *CTR) argsCheck(args *CTRArgs) error {
	if args.RK == nil || len(args.RK) == 0 {
		return fmt.Errorf("Key schedule is nil or empty, cannot perform CTR encryption")
	}

	if args.FilePath == "" {
		return fmt.Errorf("File path is empty, cannot write encrypted data")
	}

	if !state.ValidKeys[args.KeySize] {
		return fmt.Errorf("Invalid key size: %d. Valid sizes are 128, 192, or 256 bits", args.KeySize)
	}

	if args.Counter == nil {
		return fmt.Errorf("Counter is nil, cannot perform CTR encryption")
	}

	return nil
}

func (c *CTR) Encrypt(chunks [4]uint32, args ModeArgs) error {
	ctrArgs, ok := args.(*CTRArgs)
	if !ok {
		return fmt.Errorf("invalid args for CTR")
	}

	if err := c.argsCheck(ctrArgs); err != nil {
		return err
	}

	encB := [4]uint32(core.SelectEncrypt(*ctrArgs.Counter, ctrArgs.RK, ctrArgs.KeySize))

	prev := bitops.MultiXOR32(encB, chunks)

	if err := stream.WriteBinaryStream(ctrArgs.FilePath, prev); err != nil {
		return fmt.Errorf("Error writing to binary stream: %v\n", err)
	}

	ctrArgs.incrementCounter()

	return nil
}

func (c *CTR) Decrypt(chunks [4]uint32, args ModeArgs) error {
	ctrArgs, ok := args.(*CTRArgs)
	if !ok {
		return fmt.Errorf("invalid args for CTR")
	}

	if err := c.argsCheck(ctrArgs); err != nil {
		return err
	}

	encB := [4]uint32(core.SelectEncrypt(*ctrArgs.Counter, ctrArgs.RK, ctrArgs.KeySize))

	text := bitops.MultiXOR32(encB, chunks)

	if err := stream.WriteBinaryStream(ctrArgs.FilePath, text); err != nil {
		return fmt.Errorf("Error writing to binary stream: %v\n", err)
	}

	ctrArgs.incrementCounter()

	return nil
}

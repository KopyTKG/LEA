package modes

import (
	"fmt"
	"lea/core"
	"lea/state"
)

type ECB struct{}

type ECBArgs struct {
	FilePath string
	RK       []uint32
	KeySize  int
}

func (e *ECB) argsCheck(args *ECBArgs) error {
	if args.RK == nil || len(args.RK) == 0 {
		return fmt.Errorf("Key schedule is nil or empty, cannot perform ECB encryption")
	}

	if args.FilePath == "" {
		return fmt.Errorf("File path is empty, cannot write encrypted data")
	}

	if !state.ValidKeys[args.KeySize] {
		return fmt.Errorf("Invalid key size: %d. Valid sizes are 128, 192, or 256 bits", args.KeySize)
	}

	return nil
}

func (e *ECB) Encrypt(chunks [4]uint32, args ModeArgs) ([4]uint32, error) {
	ecbArgs, ok := args.(*ECBArgs)
	if !ok {
		return [4]uint32{}, fmt.Errorf("invalid args for ECB")
	}

	if err := e.argsCheck(ecbArgs); err != nil {
		return [4]uint32{}, err
	}

	encryptedBlock := core.SelectEncrypt(chunks, ecbArgs.RK, ecbArgs.KeySize)
	return encryptedBlock, nil
}

func (e *ECB) Decrypt(chunks [4]uint32, args ModeArgs) ([4]uint32, error) {
	ecbArgs, ok := args.(*ECBArgs)
	if !ok {
		return [4]uint32{}, fmt.Errorf("invalid args for ECB")
	}

	if err := e.argsCheck(ecbArgs); err != nil {
		return [4]uint32{}, err
	}

	decryptedBlock := core.SelectDecrypt(chunks, ecbArgs.RK, ecbArgs.KeySize)

	return decryptedBlock, nil
}

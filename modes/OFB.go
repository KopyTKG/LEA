package modes

import (
	"fmt"
	"lea/bitops"
	"lea/core"
	"lea/state"
	"lea/stream"
)

type OFB struct{}

type OFBArgs struct {
	FilePath string
	RK       []uint32
	KeySize  int
	Prev     *[4]uint32
}

func (o *OFB) argsCheck(args *OFBArgs) error {
	if args.RK == nil || len(args.RK) == 0 {
		return fmt.Errorf("Key schedule is nil or empty, cannot perform OFB encryption")
	}

	if args.FilePath == "" {
		return fmt.Errorf("File path is empty, cannot write encrypted data")
	}

	if !state.ValidKeys[args.KeySize] {
		return fmt.Errorf("Invalid key size: %d. Valid sizes are 128, 192, or 256 bits", args.KeySize)
	}

	if args.Prev == nil {
		return fmt.Errorf("Previous block is nil, cannot perform OFB encryption")
	}

	return nil
}

func (o *OFB) Encrypt(chunks [4]uint32, args ModeArgs) error {
	ofbArgs, ok := args.(*OFBArgs)
	if !ok {
		return fmt.Errorf("invalid args for OFB")
	}

	if err := o.argsCheck(ofbArgs); err != nil {
		return err
	}

	keystreamBlock := [4]uint32(core.SelectEncrypt(*ofbArgs.Prev, ofbArgs.RK, ofbArgs.KeySize))

	encryptedBlock := bitops.MultiXOR32(chunks, keystreamBlock)

	ofbArgs.Prev = &keystreamBlock

	if err := stream.WriteBinaryStream(ofbArgs.FilePath, encryptedBlock); err != nil {
		return fmt.Errorf("Error writing to binary stream: %v\n", err)
	}

	return nil
}

func (o *OFB) Decrypt(chunks [4]uint32, args ModeArgs) error {
	return o.Encrypt(chunks, args) // OFB is symmetric, so encryption and decryption are the same
}

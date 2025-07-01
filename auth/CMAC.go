package auth

import (
	"fmt"
	"lea/bitops"
	"lea/core"
)

type CMACargs struct {
	currentT [4]uint32
	RK       []uint32
	KeySize  int
}

type CMAC struct {
	K1 [4]uint32
	K2 [4]uint32
}

func leftShiftBlock(block [4]uint32) [4]uint32 {
	var out [4]uint32
	carry := uint32(0)
	for i := 3; i >= 0; i-- {
		nextCarry := (block[i] >> 31) & 1
		out[i] = (block[i] << 1) | carry
		carry = nextCarry
	}
	return out
}

func (c *CMAC) CreateKeys(args CMACargs) error {
	nulls := [4]uint32{0, 0, 0, 0}

	// Create L by encrypting a block of zeros
	L := core.SelectEncrypt(nulls, args.RK, args.KeySize)

	if L == [4]uint32{} {
		return fmt.Errorf("Failed to generate L in CMAC key creation")
	}

	// Left shift L to create K1
	c.K1 = leftShiftBlock(L)

	// If the most significant bit of L is set, XOR with 0x87
	if L[0]&0x80000000 != 0 {
		c.K1[3] ^= 0x87
	}

	// Left shift K1 to create K2
	c.K2 = leftShiftBlock(c.K1)

	// If the most significant bit of K1 is set, XOR with 0x87
	if c.K1[0]&0x80000000 != 0 {
		c.K2[3] ^= 0x87
	}

	return nil
}

func (c *CMAC) ComputeT(chunk [4]uint32, args *CMACargs, last bool) ([4]uint32, error) {
	if len(chunk) != 4 {
		return [4]uint32{}, fmt.Errorf("Invalid chunk size: expected 4 uint32 values, got %d", len(chunk))
	}

	if args.currentT == [4]uint32{} {
		args.currentT = [4]uint32{0, 0, 0, 0}
	}
	xorred := [4]uint32{}
	if last {
		xorred = bitops.MultiXOR32(chunk, args.currentT)
		xorred = bitops.MultiXOR32(xorred, c.K1)
	} else {

		xorred = bitops.MultiXOR32(chunk, args.currentT)
	}

	encryptedChunk := core.SelectEncrypt(xorred, args.RK, args.KeySize)
	args.currentT = encryptedChunk

	return args.currentT, nil
}

func Validate(t [4]uint32, args CMACargs) error {
	if len(t) != 4 {
		return fmt.Errorf("Invalid tag size: expected 4 uint32 values, got %d", len(t))
	}

	if args.currentT == [4]uint32{} {
		return fmt.Errorf("Current T is not set, cannot validate CMAC")
	}

	if args.currentT != t {
		return fmt.Errorf("CMAC validation failed: expected %v, got %v", args.currentT, t)
	}

	return nil
}

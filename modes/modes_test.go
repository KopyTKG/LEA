package modes

import (
	"lea/schedule"
	"testing"
)

const (
	KEYLENGTH = 256
)

var KEY [8]uint32 = [8]uint32{0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000}
var SEED [8]uint32 = [8]uint32{0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000}
var IV [4]uint32 = [4]uint32{0x00000000, 0x00000000, 0x00000000, 0x00000000}

func TestECB(t *testing.T) {
	ecb := &ECB{}

	rk, err := schedule.KeySchedule(KEYLENGTH, KEY[:], SEED[:])
	if err != nil {
		t.Fatalf("Key schedule failed: %v", err)
	}

	args := &ECBArgs{
		FilePath: "test_ecb.txt",
		RK:       rk,
		KeySize:  KEYLENGTH,
	}

	chunks := [4]uint32{0x00000001, 0x00000002, 0x00000003, 0x00000004}

	encrypted, err := ecb.Encrypt(chunks, args)
	if err != nil {
		t.Fatalf("ECB encryption failed: %v", err)
	}

	decrypted, err := ecb.Decrypt(encrypted, args)
	if err != nil {
		t.Fatalf("ECB decryption failed: %v", err)
	}

	if decrypted != chunks {
		t.Errorf("ECB decryption did not return original chunks. Got: %v, Expected: %v", decrypted, chunks)
	}
}

func TestCBC(t *testing.T) {
	cbc := &CBC{}

	ivEnc := IV
	rk, err := schedule.KeySchedule(KEYLENGTH, KEY[:], SEED[:])
	if err != nil {
		t.Fatalf("Key schedule failed: %v", err)
	}

	args := &CBCArgs{
		FilePath: "test_cbc.txt",
		RK:       rk,
		KeySize:  KEYLENGTH,
		Prev:     &ivEnc,
	}

	chunks := [4]uint32{0x00000001, 0x00000002, 0x00000003, 0x00000004}

	encrypted, err := cbc.Encrypt(chunks, args)
	if err != nil {
		t.Fatalf("CBC encryption failed: %v", err)
	}

	ivDec := IV
	argsDec := &CBCArgs{
		FilePath: "test_cbc.txt",
		RK:       rk,
		KeySize:  KEYLENGTH,
		Prev:     &ivDec,
	}

	decrypted, err := cbc.Decrypt(encrypted, argsDec)
	if err != nil {
		t.Fatalf("CBC decryption failed: %v", err)
	}

	if decrypted != chunks {
		t.Errorf("CBC decryption did not return original chunks. Got: %v, Expected: %v", decrypted, chunks)
	}
}

func TestCFB(t *testing.T) {
	cfb := &CFB{}

	ivEnc := IV
	rk, err := schedule.KeySchedule(KEYLENGTH, KEY[:], SEED[:])
	if err != nil {
		t.Fatalf("Key schedule failed: %v", err)
	}

	args := &CFBArgs{
		FilePath: "test_cfb.txt",
		RK:       rk,
		KeySize:  KEYLENGTH,
		Prev:     &ivEnc,
	}

	chunks := [4]uint32{0x00000001, 0x00000002, 0x00000003, 0x00000004}

	encrypted, err := cfb.Encrypt(chunks, args)
	if err != nil {
		t.Fatalf("CFB encryption failed: %v", err)
	}

	ivDec := IV
	argsDec := &CFBArgs{
		FilePath: "test_cfb.txt",
		RK:       rk,
		KeySize:  KEYLENGTH,
		Prev:     &ivDec,
	}
	decrypted, err := cfb.Decrypt(encrypted, argsDec)
	if err != nil {
		t.Fatalf("CFB decryption failed: %v", err)
	}

	if decrypted != chunks {
		t.Errorf("CFB decryption did not return original chunks. Got: %v, Expected: %v", decrypted, chunks)
	}
}

func TestOFB(t *testing.T) {
	ofb := &OFB{}

	rk, err := schedule.KeySchedule(KEYLENGTH, KEY[:], SEED[:])
	if err != nil {
		t.Fatalf("Key schedule failed: %v", err)
	}

	inEnc := IV
	args := &OFBArgs{
		FilePath: "test_ofb.txt",
		RK:       rk,
		KeySize:  KEYLENGTH,
		Prev:     &inEnc,
	}

	chunks := [4]uint32{0x00000001, 0x00000002, 0x00000003, 0x00000004}

	encrypted, err := ofb.Encrypt(chunks, args)
	if err != nil {
		t.Fatalf("OFB encryption failed: %v", err)
	}

	inDec := IV
	argsDec := &OFBArgs{
		FilePath: "test_ofb.txt",
		RK:       rk,
		KeySize:  KEYLENGTH,
		Prev:     &inDec,
	}
	decrypted, err := ofb.Decrypt(encrypted, argsDec)
	if err != nil {
		t.Fatalf("OFB decryption failed: %v", err)
	}

	if decrypted != chunks {
		t.Errorf("OFB decryption did not return original chunks. Got: %v, Expected: %v", decrypted, chunks)
	}
}

func TestCTR(t *testing.T) {
	ctr := &CTR{}

	rk, err := schedule.KeySchedule(KEYLENGTH, KEY[:], SEED[:])
	if err != nil {
		t.Fatalf("Key schedule failed: %v", err)
	}

	counter := [4]uint32{0x00000000, 0x00000000, 0x00000000, 0x00000000}

	args := &CTRArgs{
		FilePath: "test_ctr.txt",
		RK:       rk,
		KeySize:  KEYLENGTH,
		Counter:  &counter,
	}

	chunks := [4]uint32{0x00000001, 0x00000002, 0x00000003, 0x00000004}

	encrypted, err := ctr.Encrypt(chunks, args)
	if err != nil {
		t.Fatalf("CTR encryption failed: %v", err)
	}

	// Reset counter for decryption
	counter = [4]uint32{0x00000000, 0x00000000, 0x00000000, 0x00000000}
	args = &CTRArgs{
		FilePath: "test_ctr.txt",
		RK:       rk,
		KeySize:  KEYLENGTH,
		Counter:  &counter,
	}
	decrypted, err := ctr.Decrypt(encrypted, args)
	if err != nil {
		t.Fatalf("CTR decryption failed: %v", err)
	}

	if decrypted != chunks {
		t.Errorf("CTR decryption did not return original chunks. Got: %v, Expected: %v", decrypted, chunks)
	}
}

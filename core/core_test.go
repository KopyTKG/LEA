package core

import (
	"lea/fingerprint"
	"lea/schedule"
	"lea/types"
	"testing"
)

func TestEnc128(t *testing.T) {
	source := fingerprint.LoadSource([]byte{1, 2, 3, 4})
	key := fingerprint.Fingerprint128(source)

	rk := types.Rk128(schedule.KeySchedule(128, key[:], key[:]))

	base := [4]uint32{10,20,30,40}
	enc := Encrypt128(base, rk)
	dec := Decrypt128(enc, rk)

	for i := 0 ; i < len(base); i++ {
		if base[i] != dec[i] {
			t.Errorf("%d is not same as %d", base[i] , dec[i])
		}
	}

}

func TestEnc192(t *testing.T) {
	source := fingerprint.LoadSource([]byte{1, 2, 3, 4})
	key := fingerprint.Fingerprint192(source)

	rk := types.Rk192(schedule.KeySchedule(192, key[:], key[:]))

	base := [4]uint32{10,20,30,40}
	enc := Encrypt192(base, rk)
	dec := Decrypt192(enc, rk)

	for i := 0 ; i < len(base); i++ {
		if base[i] != dec[i] {
			t.Errorf("%d is not same as %d", base[i] , dec[i])
		}
	}

}


func TestEnc256(t *testing.T) {
	source := fingerprint.LoadSource([]byte{1, 2, 3, 4})
	key := fingerprint.Fingerprint256(source)

	rk := types.Rk256(schedule.KeySchedule(256, key[:], key[:]))

	base := [4]uint32{10,20,30,40}
	enc := Encrypt256(base, rk)
	dec := Decrypt256(enc, rk)

	for i := 0 ; i < len(base); i++ {
		if base[i] != dec[i] {
			t.Errorf("%d is not same as %d", base[i] , dec[i])
		}
	}

}

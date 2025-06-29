package modes

type ModeArgs interface{}

type Mode interface {
	Encrypt(chunks [4]uint32, args ModeArgs) error
	Decrypt(chunks [4]uint32, args ModeArgs) error
}

var SelectedMode Mode = nil
var CryptMethod func(chunks [4]uint32, args ModeArgs) error

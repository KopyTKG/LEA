package modes

type ModeArgs interface{}

type Mode interface {
	Encrypt(chunks [4]uint32, args ModeArgs) ([4]uint32, error)
	Decrypt(chunks [4]uint32, args ModeArgs) ([4]uint32, error)
}

var SelectedMode Mode = nil
var CryptMethod func(chunks [4]uint32, args ModeArgs) ([4]uint32, error)

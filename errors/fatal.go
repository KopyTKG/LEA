package errors

import "log"


func PaddingError() {
	log.Fatal("Padding overflow error")
}

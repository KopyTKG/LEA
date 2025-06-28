package help

import "fmt"

var VERSION string

func Version() {
	fmt.Printf("lea %s\n", VERSION)
}

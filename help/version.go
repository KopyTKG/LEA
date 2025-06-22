package help

import "fmt"

var VERSION string = "v2.0.1"

func Version() {
	fmt.Printf("lea %s\n", VERSION)
}

package help

import "fmt"

var VERSION string = "v3.0.0"

func Version() {
	fmt.Printf("lea %s\n", VERSION)
}

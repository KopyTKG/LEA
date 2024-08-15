package help

import "fmt"

var VERSION string = "v1.6.0"

func Version() {
    fmt.Printf("lea %s\n", VERSION)
}

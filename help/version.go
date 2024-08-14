package help

import "fmt"

var version string = "v1.6.0"

func Version() {
    fmt.Printf("lea %s\n", version)
}

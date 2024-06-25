package help

import "fmt"

var version string = "v1.1.0"

func Version() {
    fmt.Printf("lea %s\n", version)
}

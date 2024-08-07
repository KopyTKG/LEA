package help

import "fmt"

var version string = "v1.5.1"

func Version() {
    fmt.Printf("lea %s\n", version)
}

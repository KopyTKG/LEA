package help

import "fmt"

var version string = "v1.4.2"

func Version() {
    fmt.Printf("lea %s\n", version)
}

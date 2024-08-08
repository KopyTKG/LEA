package help

import "fmt"

var version string = "v1.5.4"

func Version() {
    fmt.Printf("lea %s\n", version)
}

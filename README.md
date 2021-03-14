# exit

Package exit implements an error-based alternative to os.Exit.

## Usage

```go
package main

import (
    "os"

    "github.com/azazeal/exit"
)

func main() {
    err := run()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
    }
    exit.With(err)
}
```
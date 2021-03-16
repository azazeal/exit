# exit

Package exit implements an error-based alternative to os.Exit.

## Usage

```go
package main

import (
	"io"
	"log"
	"net"
	"time"

	"github.com/azazeal/exit"
)

const (
	_       = iota + 1 // 1 is reserved for stdlib (flags, panics, etc)
	ecDial             // failed opening the connection
	ecWrite            // failed writing to the connection
	ecClose            // failed closing the connection
)

func main() {
	err := run()
	if err != nil {
		log.Println(err)
	}
	exit.With(err)
}

func run() (err error) {
	var conn net.Conn
	if conn, err = dial(); err != nil {
		return
	}

	defer func() {
		if e := close(conn); err == nil {
			err = e
		}
	}()

	err = write(conn)

	return
}

func dial() (conn net.Conn, err error) {
	const addr = "localhost:6379"

	if conn, err = net.DialTimeout("tcp", addr, time.Second<<1); err != nil {
		err = exit.Wrapf(ecDial, "failed dialing: %w", err)
	}
	return
}

func close(conn net.Conn) (err error) {
	if err = conn.Close(); err != nil {
		err = exit.Wrapf(ecClose, "failed closing: %w", err)
	}
	return
}

func write(conn net.Conn) (err error) {
	if _, err = io.WriteString(conn, "VERSION\r\n"); err != nil {
		err = exit.Wrapf(ecWrite, "failed writing: %w", err)
	}
	return
}

```

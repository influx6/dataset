package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	flag.Parse()

	command := flag.Arg(0)

	if command == "Transform" {
		var buf bytes.Buffer
		io.Copy(&buf, os.Stdin)

		fmt.Fprintf(os.Stdout, `[{"total":1, "records": %+q}]`, buf.String())
	}
}

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func main() {
	var buf bytes.Buffer
	io.Copy(&buf, os.Stdin)

	fmt.Fprintf(os.Stdout, `[{"total":1, "records": %+q}]`, buf.String())
}

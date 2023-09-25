package logo

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

var (
	Logger       io.Writer
	outputWriter *bufio.Writer
)

func init() {
	SetOutput(io.MultiWriter(os.Stdout))
}

func SetOutput(writer io.Writer) {
	outputWriter = bufio.NewWriter(writer)
}

func Printf(format string, a ...any) {
	_, _ = fmt.Fprintf(outputWriter, format, a...)
	_ = outputWriter.Flush()
}

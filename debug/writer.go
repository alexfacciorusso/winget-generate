package debug

import (
	"io"

	"github.com/fatih/color"
)

type debugWriter struct {
	writer io.Writer
}

func (w *debugWriter) Write(p []byte) (n int, err error) {
	return w.writer.Write([]byte(color.HiBlackString(string(p))))
}

// DebugWriter is a writer to be used for debugging purposes. It prints in yellow by default.
var DebugWriter = &debugWriter{writer: color.Error}

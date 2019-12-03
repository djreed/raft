package logging

import (
	"io"
	"log"
	"os"
)

const LOG_FLAGS = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile

var StdOutLogger = MakeLogger(os.Stdout)
var StdErrLogger = MakeLogger(os.Stderr)

var LOG = StdOutLogger

func MakeLogger(out io.Writer) *log.Logger {
	return log.New(out, "", LOG_FLAGS)
}

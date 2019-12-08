package logging

import (
	"io"
	"log"
	"os"
)

const LOG_FLAGS = log.Ltime | log.Lmicroseconds

var StdOutLogger = MakeLogger(os.Stdout)
var StdErrLogger = MakeLogger(os.Stderr)

var OUT = StdOutLogger
var ERR = StdErrLogger

func MakeLogger(out io.Writer) *log.Logger {
	return log.New(out, "", LOG_FLAGS)
}

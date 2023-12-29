package common

import (
	"github.com/namsral/flag"
)

var ShowLog = false

func init() {
	flag.BoolVar(&ShowLog, "show-log", false, "show log")
	flag.Parse()
}

const (
	KeyCompFiber = "fiber"
)

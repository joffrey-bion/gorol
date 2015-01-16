package log

import (
	"fmt"
	"log"
	"os"
)

var (
	logger *log.Logger = log.New(os.Stdout, "", 0)
)

const (
	lvlError int = iota
	lvlWarning
	lvlInfo
	lvlDebug
	lvlVerbose
)

const (
	LEVEL int = lvlVerbose
)

func printLog(level int, args ...interface{}) {
	if level <= LEVEL {
		logger.Println(args...)
	}
}

func E(args ...interface{}) {
	printLog(lvlError, args...)
}

func W(args ...interface{}) {
	printLog(lvlWarning, args...)
}

func I(args ...interface{}) {
	printLog(lvlInfo, args...)
}

func D(args ...interface{}) {
	printLog(lvlDebug, args...)
}

func V(args ...interface{}) {
	printLog(lvlVerbose, args...)
}

func Ef(formatStr string, args ...interface{}) {
	E(fmt.Sprintf(formatStr, args...))
}

func Wf(formatStr string, args ...interface{}) {
	W(fmt.Sprintf(formatStr, args...))
}

func If(formatStr string, args ...interface{}) {
	I(fmt.Sprintf(formatStr, args...))
}

func Df(formatStr string, args ...interface{}) {
	D(fmt.Sprintf(formatStr, args...))
}

func Vf(formatStr string, args ...interface{}) {
	V(fmt.Sprintf(formatStr, args...))
}

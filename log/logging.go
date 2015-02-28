package log

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	lvlError int = iota
	lvlWarning
	lvlInfo
	lvlDebug
	lvlVerbose
)

const (
	LEVEL  int    = lvlDebug
	INDENT string = "   "
)

var (
	logger    *log.Logger = log.New(os.Stdout, "", 0)
	indentLvl int         = 0
)

func printLog(level int, formatStr string, args ...interface{}) {
	if level <= LEVEL {
		ind := strings.Repeat(INDENT, indentLvl)
		logger.Println(fmt.Sprintf(ind+formatStr, args...))
	}
}

func Indent() {
	indentLvl++
}

func Unindent(depth int) {
	indentLvl -= depth
}

func E(formatStr string, args ...interface{}) {
	printLog(lvlError, formatStr, args...)
}

func W(formatStr string, args ...interface{}) {
	printLog(lvlWarning, formatStr, args...)
}

func I(formatStr string, args ...interface{}) {
	printLog(lvlInfo, formatStr, args...)
}

func D(formatStr string, args ...interface{}) {
	printLog(lvlDebug, formatStr, args...)
}

func V(formatStr string, args ...interface{}) {
	printLog(lvlVerbose, formatStr, args...)
}

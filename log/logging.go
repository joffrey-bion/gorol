package log

import (
	"fmt"
	"log"
	"os"
)

var (
	logger *log.Logger = log.New(os.Stdout, "bot: ", 0)
)

func E(args ...interface{}) {
	logger.Println(args)
}

func Ef(formatStr string, args ...interface{}) {
	E(fmt.Sprintf(formatStr, args))
}

func W(args ...interface{}) {
	logger.Println(args)
}

func I(args ...interface{}) {
	logger.Println(args)
}

func D(args ...interface{}) {
	logger.Println(args)
}

func V(args ...interface{}) {
	logger.Println(args)
}

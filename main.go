package main

import (
	"code.google.com/p/gcfg"
	"fmt"
	"github.com/joffrey-bion/gorol/bot"
	"github.com/joffrey-bion/gorol/config"
	"os"
)

const (
	DEFAULT_FILENAME = "default.rol"
)

func main() {
	var conf config.Config
	filename := DEFAULT_FILENAME
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}
	err := gcfg.ReadFileInto(&conf, filename)
	if err != nil {
		fmt.Printf("Config file '%s' not found.", filename)
		return
	}

	fmt.Printf("Reading config from '%s'...\n", filename)
	fmt.Println(conf)

	bot.Run(conf)
}

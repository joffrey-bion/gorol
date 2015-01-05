package bot

import (
	"fmt"
	"github.com/joffrey-bion/gorol/api"
	"github.com/joffrey-bion/gorol/config"
)

func Run(conf config.Config) bool {
	success := api.Login(conf.Account.Login, conf.Account.Password)
	if success {
		fmt.Println("Login successful as " + conf.Account.Login)
		return true
	} else {
		fmt.Println("Login failed for " + conf.Account.Login)
		return false
	}
}

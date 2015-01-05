package model

import (
	"fmt"
)

type AccountState struct {
	Gold        int
	ChestGold   int
	Mana        int
	Adventurins int
	Turns       int
}

func (s AccountState) String() string {
	return fmt.Sprintf("{gold=%d, chest=%d, turns=%d}", s.Gold, s.ChestGold, s.Turns)
}
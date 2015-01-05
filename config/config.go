package config

import (
	"fmt"
	"strings"
)

type Account struct {
	Login    string
	Password string
}

func (a Account) String() string {
	return "Account:\n   username: " + a.Login + "\n   password: " + strings.Repeat("*", len(a.Password))
}

type AttackParams struct {
	MaxTurns        int
	RepairingPeriod int
	StoringPeriod   int
}

func (ap AttackParams) String() string {
	return fmt.Sprintf("Attack params:\n   maxTurns: %d\n   repair every %d turns\n   store every %d turns", ap.MaxTurns, ap.RepairingPeriod, ap.StoringPeriod)
}

type PlayerFilter struct {
	MinRank       int
	MaxRank       int
	GoldThreshold int
}

func (pf PlayerFilter) String() string {
	return fmt.Sprintf("Player filter:\n   ranks: %d-%d\n   min gold: %d", pf.MinRank, pf.MaxRank, pf.GoldThreshold)
}

type SessionParams struct {
	NbOfAttacks        int
	HoursBetweenAttacks int64
}

func (sp SessionParams) String() string {
	return fmt.Sprintf("Session params:\n   number of attacks: %d\n   hours between attacks: %d", sp.NbOfAttacks, sp.HoursBetweenAttacks)
}

type Config struct {
	Account            Account
	Filter             PlayerFilter
	AttackParams       AttackParams
	SessionParams      SessionParams
}

func (c Config) String() string {
	return c.Account.String() + "\n" + c.Filter.String() + "\n" + c.AttackParams.String() + "\n" + c.SessionParams.String()
}

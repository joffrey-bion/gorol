package model

import (
	"fmt"
)

type Army int

const (
	UNDEFINED_ARMY Army = iota
	WARRIORS
	MAGES
	SUICIDERS
	HEALERS
)

func GetArmy(shortName string) Army {
	switch (shortName) {
	case "Chev.", "Guer.":
		return WARRIORS
	case "Sorc.":
		return MAGES
	case "Suic.":
		return SUICIDERS
	case "Sage", "Prét.":
		return HEALERS
	}
	return UNDEFINED_ARMY
}

type Alignment int

const (
	UNDEFINED_ALIGNMENT Alignment = iota
	SAINT
	CHEVALERESQUE
	ALTRUISTE
	JUSTE
	NEUTRE
	SANS_SCRUPULES
	VIL
	ABOMINABLE
	DEMONIAQUE
)

func GetAlignment(shortName string) Alignment {
	switch (shortName) {
	case "Sai.":
		return SAINT
	case "Che.":
		return CHEVALERESQUE
	case "Alt.":
		return ALTRUISTE
	case "Jus.":
		return JUSTE
	case "Neu.":
		return NEUTRE
	case "SsS.":
		return SANS_SCRUPULES
	case "Vil.":
		return VIL
	case "Abo.":
		return ABOMINABLE
	case "Dém.":
		return DEMONIAQUE
	}
	return UNDEFINED_ALIGNMENT
}

type Player struct {
	Rank      int
	Name      string
	Gold      int
	Army      Army
	Alignment Alignment
}

func (p Player) String() string {
	return fmt.Sprintf("{rank=%d, name=%s, gold=%d, army=%d, align=%d}", p.Rank, p.Name, p.Gold, p.Army, p.Alignment)
}

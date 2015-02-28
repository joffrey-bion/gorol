package bot

import (
	"github.com/joffrey-bion/gorol/api"
	"github.com/joffrey-bion/gorol/bot/timer"
	"github.com/joffrey-bion/gorol/config"
	"github.com/joffrey-bion/gorol/log"
	"github.com/joffrey-bion/gorol/model"
	"math"
	"sort"
)

func Run(conf config.Config) error {
	err := login(conf.Account.Login, conf.Account.Password)
	if err != nil {
		return err
	}
	players := attackRichest(conf.Filter, conf.AttackParams)
	//	if errs != nil {
	//		for _, e := range errs {
	//			log.E("%v", e)
	//		}
	//	}
	for _, p := range players {
		log.I("player: ", p)
	}
	return nil // no error
}

// login connects to the website with the specified credentials and wait for
// standard time.
func login(username, password string) error {
	log.D("Logging in with username %q...", username)
	err := api.Login(username, password)
	if err != nil {
		log.E("Login failure: %v", err)
		return err
	}
	log.I("Logged in as %s", username)
	log.I("")
	log.I("Faking redirection page delay... (this takes a few seconds)")
	log.I("")
	timer.WaitAfterLogin()
	return nil
}

// Logs out.
func logout() {
	log.D("Logging out...")
	err := api.Logout()
	if err != nil {
		log.E("Logout failure: %v", err)
	}
	log.I("Logout successful")
}

/**
 * Attacks the richest players matching the filter.
 *
 * @param filter
 *            the {@link PlayerFilter} to use to choose players
 * @param params
 *            the {@link AttackParams} to use for attacks/actions sequencing
 * @return the total gold stolen
 */
func attackRichest(filter config.PlayerFilter, params config.AttackParams) []*model.Player {
	log.I("Starting massive attack on players ranked %d to %d richer than %d gold (%d attacks max)", filter.MinRank, filter.MaxRank,
		filter.GoldThreshold, params.MaxTurns)
	log.I("Searching players matching the config filter...")
	log.Indent()
	maxPlayersToAttack := int(math.Max(float64(params.MaxTurns), float64(api.GetCurrentState().Turns)))
	matchingPlayers := model.PlayerList{}
	startRank := filter.MinRank
	for startRank < filter.MaxRank {
		log.D("Reading page of players ranked %d to %d...", startRank, startRank+98)
		playersInPage, errors := api.ListPlayers(startRank)
		if len(errors) > 0 {
			logErrors(errors)
			return []*model.Player{}
		}
		filteredPage := filter.Apply(playersInPage) // filter players
		sort.Sort(model.PlayerList(filteredPage))   // richest first
		if len(filteredPage) > maxPlayersToAttack {
			filteredPage = filteredPage[:maxPlayersToAttack] // limit number of players
		}
		maxRankToAttack := int(math.Min(float64(filter.MaxRank), float64(startRank+98)))
		log.I("%2d matching player(s) ranked %d to %d", len(filteredPage), startRank, maxRankToAttack)
		matchingPlayers = append(matchingPlayers, filteredPage...)
		timer.ReadPage()
		startRank += 99
	}
	log.Unindent(1)
	log.I("")
	nbMatchingPlayers := len(matchingPlayers)
	if nbMatchingPlayers > maxPlayersToAttack {
		log.I("%d players matching rank and gold criterias, filtering only the richest of them...", nbMatchingPlayers)
		// too many players, select only the richest
		filteredPlayers := filter.Apply(matchingPlayers) // filter players
		sort.Sort(model.PlayerList(filteredPlayers))     // richest first
		if len(filteredPlayers) > maxPlayersToAttack {
			filteredPlayers = filteredPlayers[:maxPlayersToAttack] // limit number of players
		}
		//return attackAll(filteredPlayers, params)
		return filteredPlayers
	} else {
		//return attackAll(matchingPlayers, params)
		return matchingPlayers
	}
}

/**
 * Attacks all the specified players, following the given parameters.
 *
 * @param playersToAttack
 *            the filtered list of players to attack. They must verify the thresholds specified
 *            by the given {@link PlayerFilter}.
 * @param params
 *            the parameters to follow. In particular the storing and repair frequencies are
 *            used.
 * @return the total gold stolen
 */
func attackAll(playersToAttack []*model.Player, params config.AttackParams) int {
	log.I("%d players to attack", len(playersToAttack))
	if api.GetCurrentState().Turns == 0 {
		log.E("No turns available, impossible to attack.")
		return 0
	} else if api.GetCurrentState().Turns < len(playersToAttack) {
		log.E("Not enough turns to attack this many players, attack aborted.")
		return 0
	}
	totalGoldStolen := 0
	nbConsideredPlayers := 0
	nbAttackedPlayers := 0
	for _, player := range playersToAttack {
		nbConsideredPlayers++
		// attack player
		goldStolen := attack(player)
		if goldStolen < 0 {
			// error, player not attacked
			continue
		}
		totalGoldStolen += goldStolen
		nbAttackedPlayers++
		isLastPlayer := nbConsideredPlayers == len(playersToAttack)
		// repair weapons as specified
		if nbAttackedPlayers%params.RepairingPeriod == 0 || isLastPlayer {
			timer.ChangePage()
			repairWeapons()
		}
		// store gold as specified
		if nbAttackedPlayers%params.StoringPeriod == 0 || isLastPlayer {
			timer.ChangePage()
			storeGoldIntoChest()
			timer.PauseWhenSafe()
		} else {
			timer.ChangePageLong()
		}
	}
	log.I("%d total gold stolen from %d players", totalGoldStolen, nbAttackedPlayers)
	log.I("The chest now contains %d gold.", api.GetCurrentState().ChestGold)
	return totalGoldStolen
}

/**
 * Attacks the specified player.
 *
 * @param player
 *            the player to attack
 * @return the gold stolen from that player
 */
func attack(player *model.Player) int {
	log.D("Attacking player %s...", player.Name)
	log.Indent()
	log.V("Displaying player page...")
	playerGold, err := api.DisplayPlayer(player.Name)
	log.Indent()
	if err != nil {
		log.E("Something's wrong: request failed")
		log.Unindent(2)
		return -1
	} else if playerGold != player.Gold {
		log.W("Something's wrong: the player does not have the expected gold")
		log.Unindent(2)
		return -1
	}
	log.Unindent(1)

	timer.ActionInPage()

	log.V("Attacking...")
	goldStolen, errors := api.Attack(player.Name)
	log.Unindent(1)
	if len(errors) > 0 {
		logErrors(errors)
		return 0
	}
	if goldStolen > 0 {
		log.I("Victory! %d gold stolen from player %s, current gold: %d", goldStolen, player.Name, api.GetCurrentState().Gold)
	} else {
		log.W("Defeat! Ach, player %s was too sronk! Current gold: %s", player.Name,
			api.GetCurrentState().Gold)
	}
	return goldStolen
}

func storeGoldIntoChest() int {
	log.V("Storing gold into the chest...")
	log.Indent()
	log.V("Displaying chest page...")
	log.Indent()
	amount, err := api.DisplayChestPage()
	if err != nil {
		log.E("Something went wrong when displaying the page")
		log.Unindent(2)
		return 0
	}
	log.Unindent(1)
	log.V("%d gold to store", amount)

	timer.ActionInPage()

	log.V("Storing everything...")
	log.Indent()
	err = api.StoreInChest(amount)
	if err != nil {
		log.V("The gold is safe!")
	} else {
		log.E("Something went wrong!")
	}
	log.Unindent(2)
	log.I("%d gold stored in chest, total: %d", amount, api.GetCurrentState().ChestGold)
	return amount
}

func repairWeapons() {
	log.V("Repairing weapons...")
	log.Indent()
	log.V("Displaying weapons page...")
	wornness, err := api.DisplayWeaponsPage()
	log.V("Weapons worn at %d%%", wornness)

	timer.ActionInPage()

	log.V("Repair request...")
	err = api.RepairWeapons()
	log.Unindent(1)
	if err != nil {
		log.E("Couldn't repair weapons, is there enough gold?")
	} else {
		log.I("Weapons repaired")
	}
}

func logErrors(errors []error) {
	for _, e := range errors {
		log.E("%v", e)
	}
}

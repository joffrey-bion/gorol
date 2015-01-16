package bot

import (
	"github.com/joffrey-bion/gorol/api"
	"github.com/joffrey-bion/gorol/config"
	"github.com/joffrey-bion/gorol/log"
)

func Run(conf config.Config) error {
	err := login(conf.Account.Login, conf.Account.Password)
	if err != nil {
		return err
	}
	players, errs := api.ListPlayers(1800)
	if errs != nil {
		for _, e := range errs {
			log.E(e);
		}
	}
	for _, p := range players {
		log.I("player: ", p)
	}
	return nil // no error
}

// login connects to the website with the specified credentials and wait for
// standard time.
func login(username, password string) error {
	log.D("Logging in with username ", username, "...")
	err := api.Login(username, password)
	if err != nil {
		log.E("Login failure: ", err)
		return err
	}
	log.I("Logged in with username: ", username)
	log.I("")
	log.I("Faking redirection page delay... (this takes a few seconds)")
	log.I("")
	WaitAfterLogin()
	return nil
}

// Logs out.
func logout() {
	log.D("Logging out...")
	err := api.Logout()
	if err != nil {
		log.E("Logout failure: ", err)
	}
	log.I("Logout successful")
}

//    /**
//     * Attacks the richest players matching the filter.
//     *
//     * @param filter
//     *            the {@link PlayerFilter} to use to choose players
//     * @param params
//     *            the {@link AttackParams} to use for attacks/actions sequencing
//     * @return the total gold stolen
//     */
//    private int attackRichest(PlayerFilter filter, AttackParams params) {
//        log.I("Starting massive attack on players ranked ", filter.getMinRank(), " to ", filter.getMaxRank(),
//                " richer than ", filter.getGoldThreshold(), " gold (", params.getMaxTurns(), " attacks max)")
//        log.I("Searching players matching the config filter...")
//        log.indent()
//        List<Player> matchingPlayers = new ArrayList<>()
//        int startRank = filter.getMinRank()
//        while (startRank < filter.getMaxRank()) {
//            log.D("Reading page of players ranked ", startRank, " to ", startRank + 98, "...")
//            List<Player> filteredPage = api.listPlayers(startRank).stream() // stream players
//                    .filter(p -> p.getGold() >= filter.getGoldThreshold()) // above gold threshold
//                    .filter(p -> p.getRank() <= filter.getMaxRank()) // below max rank
//                    .sorted(richestFirst) // richest first
//                    .limit(params.getMaxTurns()) // limit to max turns
//                    .limit(api.getCurrentState().turns) // limit to available turns
//                    .collect(Collectors.toList())
//            log.I(TAG,
//                    String.format("%2d matching player(s) ranked %d to %d", filteredPage.size(), startRank,
//                            Math.min(startRank + 98, filter.getMaxRank())))
//            matchingPlayers.addAll(filteredPage)
//            fakeTime.readPage()
//            startRank += 99
//        }
//        log.deindent(1)
//        log.I("")
//        int nbMatchingPlayers = matchingPlayers.size()
//        if (nbMatchingPlayers > params.getMaxTurns() || nbMatchingPlayers > api.getCurrentState().turns) {
//            log.I(matchingPlayers.size(),
//                    " players matching rank and gold criterias, filtering only the richest of them...")
//            // too many players, select only the richest
//            List<Player> playersToAttack = matchingPlayers.stream() // stream players
//                    .sorted(richestFirst) // richest first
//                    .limit(params.getMaxTurns()) // limit to max turns
//                    .limit(api.getCurrentState().turns) // limit to available turns
//                    .collect(Collectors.toList())
//            return attackAll(playersToAttack, params)
//        } else {
//            return attackAll(matchingPlayers, params)
//        }
//    }
//
//    /**
//     * Attacks all the specified players, following the given parameters.
//     *
//     * @param playersToAttack
//     *            the filtered list of players to attack. They must verify the thresholds specified
//     *            by the given {@link PlayerFilter}.
//     * @param params
//     *            the parameters to follow. In particular the storing and repair frequencies are
//     *            used.
//     * @return the total gold stolen
//     */
//    private int attackAll(List<Player> playersToAttack, AttackParams params) {
//        log.I(playersToAttack.size(), " players to attack")
//        if (api.getCurrentState().turns == 0) {
//            log.E("No turns available, impossible to attack.")
//            return 0
//        } else if (api.getCurrentState().turns < playersToAttack.size()) {
//            log.E("Not enough turns to attack this many players, attack aborted.")
//            return 0
//        }
//        int totalGoldStolen = 0
//        int nbConsideredPlayers = 0
//        int nbAttackedPlayers = 0
//        for (Player player : playersToAttack) {
//            nbConsideredPlayers++
//            // attack player
//            int goldStolen = attack(player)
//            if (goldStolen < 0) {
//                // error, player not attacked
//                continue
//            }
//            totalGoldStolen += goldStolen
//            nbAttackedPlayers++
//            boolean isLastPlayer = nbConsideredPlayers == playersToAttack.size()
//            // repair weapons as specified
//            if (nbAttackedPlayers % params.getRepairFrequency() == 0 || isLastPlayer) {
//                fakeTime.changePage()
//                repairWeapons()
//            }
//            // store gold as specified
//            if (nbAttackedPlayers % params.getStoringFrequency() == 0 || isLastPlayer) {
//                fakeTime.changePage()
//                storeGoldIntoChest()
//                fakeTime.pauseWhenSafe()
//            } else {
//                fakeTime.changePageLong()
//            }
//        }
//        log.I(totalGoldStolen, " total gold stolen from ", nbAttackedPlayers, " players")
//        log.I("The chest now contains ", api.getCurrentState().chestGold, " gold.")
//        return totalGoldStolen
//    }
//
//    /**
//     * Attacks the specified player.
//     *
//     * @param player
//     *            the player to attack
//     * @return the gold stolen from that player
//     */
//    private int attack(Player player) {
//        log.D("Attacking player ", player.getName(), "...")
//        log.indent()
//        log.V("Displaying player page...")
//        int playerGold = api.displayPlayer(player.getName())
//        log.indent()
//        if (playerGold == RoLAdapter.ERROR_REQUEST) {
//            log.E("Something's wrong: request failed")
//            log.deindent(2)
//            return -1
//        } else if (playerGold != player.getGold()) {
//            log.W("Something's wrong: the player does not have the expected gold")
//            log.deindent(2)
//            return -1
//        }
//        log.deindent(1)
//
//        fakeTime.actionInPage()
//
//        log.V("Attacking...")
//        int goldStolen = api.attack(player.getName())
//        log.deindent(1)
//        if (goldStolen > 0) {
//            log.I("Victory! ", goldStolen, " gold stolen from player ", player.getName(), ", current gold: ",
//                    api.getCurrentState().gold)
//        } else if (goldStolen == RoLAdapter.ERROR_STORM_ACTIVE) {
//            log.E("Cannot attack: a storm is raging upon your kingdom!")
//        } else if (goldStolen == RoLAdapter.ERROR_REQUEST) {
//            log.E("Attack request failed, something went wrong")
//        } else {
//            log.W("Defeat! Ach, player ", player.getName(), " was too sronk! Current gold: ",
//                    api.getCurrentState().gold)
//        }
//        return goldStolen
//    }
//
//    private int storeGoldIntoChest() {
//        log.V("Storing gold into the chest...")
//        log.indent()
//        log.V("Displaying chest page...")
//        int amount = api.displayChestPage()
//        log.V(amount + " gold to store")
//
//        fakeTime.actionInPage()
//
//        log.V("Storing everything...")
//        log.indent()
//        boolean success = api.storeInChest(amount)
//        if (success) {
//            log.V("The gold is safe!")
//        } else {
//            log.V("Something went wrong!")
//        }
//        log.deindent(2)
//        log.I(amount, " gold stored in chest, total: " + api.getCurrentState().chestGold)
//        return amount
//    }
//
//    private void repairWeapons() {
//        log.V("Repairing weapons...")
//        log.indent()
//        log.V("Displaying weapons page...")
//        int wornness = api.displayWeaponsPage()
//        log.V("Weapons worn at ", wornness, "%")
//
//        fakeTime.actionInPage()
//
//        log.V("Repair request...")
//        boolean success = api.repairWeapons()
//        log.deindent(1)
//        if (!success) {
//            log.E("Couldn't repair weapons, is there enough gold?")
//        } else {
//            log.I("Weapons repaired")
//        }
//    }

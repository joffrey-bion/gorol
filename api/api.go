package api

import (
	"fmt"
	"github.com/joffrey-bion/gorol/api/net"
	"github.com/joffrey-bion/gorol/api/parser"
	"github.com/joffrey-bion/gorol/model"
	"math/rand"
	"net/url"
	"strconv"
	"time"
)

const (
	BASE_URL  string = "http://www.riseoflords.com/"
	URL_INDEX string = BASE_URL + "index.php"
	URL_GAME  string = BASE_URL + "jeu.php"

	PAGE_LOGIN        string = "verifpass"
	PAGE_LOGOUT       string = "logout"
	PAGE_USERS_LIST   string = "main/conseil_de_guerre"
	PAGE_USER_DETAILS string = "main/fiche"
	PAGE_ATTACK       string = "main/combats"
	PAGE_CHEST        string = "main/tresor"
	PAGE_WEAPONS      string = "main/arsenal"
	PAGE_SORCERY      string = "main/autel_sorciers"

	ERROR_REQUEST      int = -1
	ERROR_STORM_ACTIVE int = -2
)

var (
	state model.AccountState
)

func init() {
	rand.Seed(time.Now().UnixNano())
	parser.SetBaseURL(BASE_URL)
}

func GetCurrentState() model.AccountState {
	return state
}

func randomCoord(min int, max int) string {
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

// Login performs the login request with the specified credentials. One needs to
// wait at least 5-6 seconds to fake real login delay after success.
func Login(username string, password string) error {
	form := url.Values{}
	form.Add("LogPseudo", username)
	form.Add("LogPassword", password)
	respBody, err := net.Post(URL_INDEX, PAGE_LOGIN, form)
	if err != nil {
		return err
	}
	if parser.Contains(respBody, "Identification incorrecte") {
		return fmt.Errorf("Login failed: wrong credentials for login %q", username)
	}
	if !parser.Contains(respBody, "Identification réussie!") {
		return fmt.Errorf("Login failed: response: %v", respBody)
	}
	return nil
}

// Logout logs the current user out.
func Logout() error {
	respBody, err := net.Get(URL_INDEX, PAGE_LOGOUT, nil)
	if err != nil {
		return err
	}
	if !parser.Contains(respBody, "Déjà inscrit? Connectez-vous") {
		return fmt.Errorf("something went wrong while logging out")
	}
	return nil
}

// Returns a list of 99 users, starting at the specified rank.
func ListPlayers(startRank int) ([]*model.Player, []error) {
	query := url.Values{}
	query.Add("Debut", strconv.Itoa(startRank+1))
	if rand.Intn(5) == 0 {
		query.Add("x", randomCoord(5, 35))
		query.Add("y", randomCoord(5, 25))
	}
	respBody, err := net.Get(URL_GAME, PAGE_USERS_LIST, query)
	if err != nil {
		return nil, []error{err}
	}
	if !parser.Contains(respBody, "Recherche pseudo:") {
		return nil, []error{fmt.Errorf("ListPlayers(%d): the page does not seem right", startRank)}
	}
	players, errors := parser.ParsePlayerList(respBody)
	updateErrors := parser.UpdateState(&state, respBody)
	return players, append(errors, updateErrors...)
}

// DisplayPlayer requests the specified player's detail page. Use this to fake a
// visit on the user detail page prior to attacking. Returns the gold of the given
// player.
func DisplayPlayer(playerName string) (int, []error) {
	query := url.Values{}
	query.Add("voirpseudo", playerName)
	respBody, err := net.Get(URL_GAME, PAGE_USER_DETAILS, query)
	if err != nil {
		return 0, []error{err}
	}
	if !parser.Contains(respBody, "Seigneur "+playerName) {
		return 0, []error{fmt.Errorf("DisplayPlayer(%s): the page does not seem right", playerName)}
	}
	gold, err := parser.ParsePlayerGold(respBody)
	updateErrors := parser.UpdateState(&state, respBody)
	return gold, append([]error{err}, updateErrors...)
}

// Attacks the specified user with one game turn. Returns the gold stolen during the
// attack.
func Attack(username string) (int, []error) {
	form := url.Values{}
	form.Add("a", "ok")
	form.Add("PseudoDefenseur", username)
	form.Add("NbToursToUse", "1")
	respBody, err := net.Post(URL_GAME, PAGE_ATTACK, form)
	if err != nil {
		return 0, []error{err}
	}
	if parser.Contains(respBody, "remporte le combat!") {
		gold, err := parser.ParseGoldStolen(respBody)
		updateErrors := parser.UpdateState(&state, respBody)
		return gold, append([]error{err}, updateErrors...)
	} else if parser.Contains(respBody, "perd cette bataille!") {
		errors := []error{fmt.Errorf("Attack(%s): defeat.", username)}
		updateErrors := parser.UpdateState(&state, respBody)
		return 0, append(errors, updateErrors...)
	} else if parser.Contains(respBody, "tempête magique s'abat") {
		errors := []error{fmt.Errorf("Attack(%s): cannot attack: a storm is raging here", username)}
		updateErrors := parser.UpdateState(&state, respBody)
		return 0, append(errors, updateErrors...)
	} else {
		return 0, []error{fmt.Errorf("Attack(%s): something went wrong, the page does not seem right", username)}
	}
}

// Displays the chest page, and returns the amount of gold that could be stored in
// the chest.
func DisplayChestPage() (int, []error) {
	respBody, err := net.Get(URL_GAME, PAGE_CHEST, url.Values{})
	if err != nil {
		return 0, []error{err}
	}
	if !parser.Contains(respBody, "ArgentAPlacer") {
		return 0, []error{fmt.Errorf("DisplayChestPage(): the page does not seem right")}
	}
	updateErrors := parser.UpdateState(&state, respBody)
	return state.Gold, updateErrors
}

// Stores the specified amount of gold into the chest. The amount has to match the current gold
// of the user, which should first be retrieved by calling DisplayChestPage().
func StoreInChest(amount int) []error {
	form := url.Values{}
	form.Add("ArgentAPlacer", strconv.Itoa(amount))
	form.Add("x", randomCoord(10, 60))
	form.Add("y", randomCoord(10, 60))
	respBody, err := net.Post(URL_GAME, PAGE_CHEST, form)
	if err != nil {
		return []error{err}
	}
	if state.Gold != 0 {
		return []error{fmt.Errorf("StoreInChest(%d): something went wrong, %d gold remaining", amount, state.Gold)}
	}
	return parser.UpdateState(&state, respBody)
}

// Displays the weapons page. Used to fake a visit on the weapons page before repairing or
// buying weapons and equipment. Returns the percentage of wornness of the weapons.
func DisplayWeaponsPage() (int, []error) {
	respBody, err := net.Get(URL_GAME, PAGE_WEAPONS, nil)
	if err != nil {
		return 0, []error{err}
	}
	if !parser.Contains(respBody, "Faites votre choix") {
		return 0, []error{fmt.Errorf("DisplayWeaponsPage(): the page does not seem right")}
	}
	wornness, err := parser.ParseWeaponsWornness(respBody)
	if err != nil {
		return wornness, []error{err}
	}
	return wornness, parser.UpdateState(&state, respBody)
}

// Repairs weapons.
func RepairWeapons() []error {
	query := url.Values{}
	query.Add("a", "repair")
	query.Add("onglet", "")
	respBody, err := net.Post(URL_GAME, PAGE_WEAPONS, query)
	if err != nil {
		return []error{err}
	}
	if !parser.Contains(respBody, "Faites votre choix") {
		return []error{fmt.Errorf("RepairWeapons(): the page does not seem right")}
	}
	return parser.UpdateState(&state, respBody)
}

//    /**
//     * Displays the sorcery page. Used to fake a visit on the sorcery page before casting a spell.
//     *
//     * @return the available mana, or {@link #ERROR_REQUEST} if the request failed
//     */
//    func DisplaySorceryPage() int {
//        HttpUriRequest request = Request.from(URL_GAME, PAGE_SORCERY).get()
//        response string = http.execute(request)
//        if (response.contains("Niveau de vos sorciers")) {
//            return state.mana
//        } else {
//            return ERROR_REQUEST
//        }
//    }
//
//    /**
//     * Casts the dissipation spell to get rid of the protective aura. Useful before self-casting a
//     * storm.
//     *
//     * @return true if the request succeeded, false otherwise
//     */
//    func DissipateProtectiveAura() bool {
//        HttpGet request = Request.from(URL_GAME, PAGE_SORCERY) //
//                .addParameter("a", "lancer") //
//                .addParameter("idsort", "14") //
//                .get()
//        response string = http.execute(request)
//        Parser.updateState(state, response)
//        return true // TODO handle failure
//    }
//
//    /**
//     * Casts a magic storm on the specified player.
//     *
//     * @param playerName
//     *            the amount of gold to store into the chest
//     * @return true if the request succeeded, false otherwise
//     */
//    func CastMagicStorm( playerName string) bool {
//        HttpPost request = Request.from(URL_GAME, PAGE_SORCERY) //
//                .addParameter("a", "lancer") //
//                .addParameter("idsort", "5") //
//                .addPostData("tempete_pseudo_cible", playerName) //
//                .post()
//        response string = http.execute(request)
//        Parser.updateState(state, response)
//        return true // TODO handle failure
//    }

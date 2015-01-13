package api

import (
	"fmt"
	"github.com/joffrey-bion/gorol/model"
	_ "io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	BASE_URL string = "http://www.riseoflords.com"
	URL_INDEX string = BASE_URL + "/index.php"
	URL_GAME  string = BASE_URL + "/jeu.php"

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
	state  model.AccountState
	logger *log.Logger = log.New(os.Stderr, "api: ", 0)
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GetCurrentState() model.AccountState {
	return state
}

func randomCoord(min int, max int) string {
	return string(rand.Intn(max-min+1) + min)
}

func gamePageUrl(page string, query url.Values) string {
	return pageUrl(URL_GAME, page) + "+" + query.Encode()
}

func pageUrl(base string, page string) string {
	return base + "?p=" + page
}

// Login performs the login request with the specified credentials.
// One needs to wait at least 5-6 seconds to fake real login delay after sucess.
// Returns true for success, false for failure.
func Login(username string, password string) error {
	form := url.Values{}
	form.Add("LogPseudo", username)
	form.Add("LogPassword", password)
	resp, err := http.PostForm(pageUrl(URL_INDEX, PAGE_LOGIN), form)
	if err != nil {
		return err
	}
	respBody, err := String(resp)
	if err != nil {
		return err
	}
	if Contains(respBody, "Identification réussie!") {
		return nil
	}
	return fmt.Errorf("something went wrong while logging in")
}

// Logout logs the current user out.
// Returns true if the request succeeded, false otherwise
func Logout() error {
	resp, err := http.Get(pageUrl(URL_INDEX, PAGE_LOGOUT))
	if err != nil {
		return err
	}
	respBody, err := String(resp)
	if err != nil {
		return err
	}
	if Contains(respBody, "Déjà inscrit? Connectez-vous") {
		return nil
	}
	return fmt.Errorf("something went wrong while logging out")
}

// Returns a list of 99 users, starting at the specified rank.
func ListPlayers(startRank int) ([]model.Player, error) {
	query := url.Values{}
	query.Add("Debut", string(startRank + 1))
    if (rand.Intn(5) == 0) {
        query.Add("x", randomCoord(5, 35))
        query.Add("y", randomCoord(5, 25))
    }
    resp, err := http.Get(gamePageUrl(PAGE_USERS_LIST, query))
    if (err != nil) {
    	return nil, err
    }
	respBody, err := String(resp)
	if err != nil {
		return nil, err
	}
    if (Contains(respBody, "Recherche pseudo:")) {
        UpdateState(&state, respBody)
        return ParsePlayerList(respBody)
    } else {
        return nil, fmt.Errorf("ListPlayers(%d): the page does not seem right", startRank)
    }
}

//    /**
//     * Displays the specified player's detail page. Used to fake a visit on the user detail page
//     * before an attack. The result does not matter.
//     * 
//     * @param playerName
//     *            the name of the player to lookup
//     * @return the specified player's current gold, or {@link #ERROR_REQUEST} if the request failed
//     */
//    func DisplayPlayer( playerName string) int {
//         request HttpGet = Request.from(URL_GAME, PAGE_USER_DETAILS) //
//                .addParameter("voirpseudo", playerName) //
//                .get()
//         response string = http.execute(request)
//        if (response.contains("Seigneur " + playerName)) {
//            Parser.updateState(state, response)
//            return Parser.parsePlayerGold(response)
//        } else {
//            return ERROR_REQUEST
//        }
//    }
//
//    /**
//     * Attacks the specified user with one game turn.
//     * 
//     * @param username
//     *            the name of the user to attack
//     * @return the gold stolen during the attack, or {@link #ERROR_REQUEST} if the request failed
//     */
//    func Attack( username string ) int {
//         request HttpPost = Request.from(URL_GAME, PAGE_ATTACK) //
//                .addParameter("a", "ok") //
//                .addPostData("PseudoDefenseur", username) //
//                .addPostData("NbToursToUse", "1") //
//                .post()
//        response string = http.execute(request)
//        if (response.contains("remporte le combat!") || response.contains("perd cette bataille!")) {
//            Parser.updateState(state, response)
//            return Parser.parseGoldStolen(response)
//        } else if (response.contains("tempête magique s'abat")) {
//            Parser.updateState(state, response)
//            return ERROR_STORM_ACTIVE
//        } else {
//            return ERROR_REQUEST
//        }
//    }
//
//    /**
//     * Gets the chest page from the server, and returns the amount of money that could be stored in
//     * the chest.
//     * 
//     * @return the amount of money that could be stored in the chest, which is the current amount of
//     *         gold of the player, or {@link #ERROR_REQUEST} if the request failed
//     */
//    func DisplayChestPage() int {
//        HttpGet request = Request.from(URL_GAME, PAGE_CHEST).get()
//         response string = http.execute(request)
//        if (response.contains("ArgentAPlacer")) {
//            Parser.updateState(state, response)
//            return state.gold
//        } else {
//            return ERROR_REQUEST
//        }
//    }
//
//    /**
//     * Stores the specified amount of gold into the chest. The amount has to match the current gold
//     * of the user, which should first be retrieved by calling {@link #displayChestPage()}.
//     * 
//     * @param amount
//     *            the amount of gold to store into the chest
//     * @return true if the request succeeded, false otherwise
//     */
//    func StoreInChest( amount int) bool {
//        HttpPost request = Request.from(URL_GAME, PAGE_CHEST) //
//                .addPostData("ArgentAPlacer", string.valueOf(amount)) //
//                .addPostData("x", randomCoord(10, 60)) //
//                .addPostData("y", randomCoord(10, 60)) //
//                .post()
//         response string = http.execute(request)
//        Parser.updateState(state, response)
//        return state.gold == 0
//    }
//
//    /**
//     * Displays the weapons page. Used to fake a visit on the weapons page before repairing or
//     * buying weapons and equipment.
//     * 
//     * @return the percentage of wornness of the weapons, or {@link #ERROR_REQUEST} if the request failed
//     */
//    func DisplayWeaponsPage() int {
//        HttpGet request = Request.from(URL_GAME, PAGE_WEAPONS).get()
//         response string = http.execute(request)
//        if (response.contains("Faites votre choix")) {
//            return Parser.parseWeaponsWornness(response)
//        } else {
//            return ERROR_REQUEST
//        }
//    }
//
//    /**
//     * Repairs weapons.
//     * 
//     * @return true if the repair succeeded, false otherwise
//     */
//    func RepairWeapons() bool {
//        HttpGet request = Request.from(URL_GAME, PAGE_WEAPONS) //
//                .addParameter("a", "repair") //
//                .addParameter("onglet", "") //
//                .get()
//         response string = http.execute(request)
//        if (response.contains("Faites votre choix")) {
//            return Parser.parseWeaponsWornness(response) == 0
//        } else {
//            return false
//        }
//    }
//
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
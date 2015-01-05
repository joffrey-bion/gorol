package api

import (
	"github.com/joffrey-bion/gorol/api/parser"
	"github.com/joffrey-bion/gorol/model"
	_ "io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
)

const (
	URL_INDEX string = "http://www.riseoflords.com/index.php"
	URL_GAME  string = "http://www.riseoflords.com/jeu.php"

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
	logger *log.Logger = log.New(os.Stderr, "[API] ", 0)
)

func GetCurrentState() model.AccountState {
	return state
}

func randomCoord(min int, max int) string {
	return string(rand.Intn(max-min+1) + min)
}

func gamePageUrl(page string) string {
	return pageUrl(URL_GAME, page)
}

func pageUrl(base string, page string) string {
	return base + "?p=" + page
}

// Login performs the login request with the specified credentials.
// One needs to wait at least 5-6 seconds to fake real login delay after sucess.
// Returns true for success, false for failure.
func Login(username string, password string) bool {
	form := url.Values{}
	form.Add("LogPseudo", username)
	form.Add("LogPassword", password)
	resp, err := http.PostForm(pageUrl(URL_INDEX, PAGE_LOGIN), form)
	if err != nil {
		logger.Println(err)
		return false
	}
	return parser.Contains(resp, "Identification réussie!")
}

// Logout logs the current user out.
// Returns true if the request succeeded, false otherwise
func Logout() bool {
	resp, err := http.Get(pageUrl(URL_INDEX, PAGE_LOGOUT))
	if err != nil {
		logger.Println(err)
		return false
	}
	return parser.Contains(resp, "Déjà inscrit? Connectez-vous")
}

// Returns a list of 99 users, starting at the specified rank.
//func ListPlayers( startRank int) []Player {
//	query := url.Values{}
//	query.Add("Debut", string(startRank + 1))
//    if (rand.Booln) {
//        query.Add("x", randomCoord(5, 35))
//        query.Add("y", randomCoord(5, 25))
//    }
//    resp, err := http.Get(gamePageUrl(PAGE_USERS_LIST) + query.Encode())
//    if (Contains(resp, err, "Recherche pseudo:")) {
//        Parser.updateState(state, response)
//        return Parser.parsePlayerList(response)
//    } else {
//        return new ArrayList<>()
//    }
//}

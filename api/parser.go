package api

import (
	"code.google.com/p/go-charset/charset"
	_ "code.google.com/p/go-charset/data"
	"fmt"
	"github.com/joffrey-bion/gorol/api/ocr"
	"github.com/joffrey-bion/gorol/model"
	"github.com/joffrey-bion/gosoup"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// assertIsTag panics if the node is not the specified tag
func assertIsTag(node *gosoup.Node, expectedTag string) {
	if !node.IsTag(expectedTag) {
		panic(fmt.Sprintf("expected tag <%s>, got <%s>", expectedTag, node.Data))
	}
}

// String returns the body of the given HTTP response as a string.
// After a call to this function, the response reader is closed and can't be used anymore.
func String(resp *http.Response) (string, error) {
	defer resp.Body.Close()
	// first parsing even in wrong charset, to get the encoding from the HTML
	doc, err := gosoup.Parse(resp.Body)
	if err != nil {
		logger.Println(err)
		return "", fmt.Errorf("String: error reading the response as HTML: %v", err)
	}
	// get the charset from the incorrect html (hopefully not that incorrect)
	cset, err := gosoup.GetDocCharset(doc)
	if err != nil {
		logger.Println(err)
		return "", fmt.Errorf("String: error getting the charset: %v", err)
	}
	// pipe to re-read the DOM tree
	preader, pwriter := io.Pipe()
	defer preader.Close()
	// rewrite the DOM tree as-is in the pipe
	go func() {
		gosoup.Render(pwriter, doc)
		pwriter.Close()
	}()
	// translate the data from the pipe to UTF-8 using the proper encoding
	reader, err := charset.NewReader(cset, preader)
	if err != nil {
		logger.Println(err)
		return "", fmt.Errorf("String: error reading the response with the appropriate charset")
	}
	// read the UTF-8 data as a String
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		logger.Println(err)
		return "", fmt.Errorf("String: error reading the response: %v", err)
	}

	return string(body), nil
}

func Contains(resp string, s string) bool {
	return strings.Contains(resp, s)
}

// Updates the specified account state based on the top elements of the specified page.
func UpdateState(state *model.AccountState, resp string) error {
	doc, err := gosoup.Parse(strings.NewReader(resp))
	if err != nil {
		return err
	}
	state.Gold = findNumValueInImgUrl(doc, "onmouseover", "A chaque tour de jeu")
	state.ChestGold = findNumValueInImgUrl(doc, "onmouseover", "Votre coffre magique")
	state.Mana = findNumValueInImgUrl(doc, "onmouseover", "Votre mana représente")
	state.Turns = findNumValueInImgUrl(doc, "onmouseover", "Un nouveau tour de jeu")
	state.Adventurins = findNumValueInImgUrl(doc, "href", "main/aventurines_detail")
	return nil
}

func findNumValueInImgUrl(node *gosoup.Node, attrKey, attrValuePart string) int {
	attrs := node.ChildrenByAttrValueContaining(attrKey, attrValuePart).First().Attrs
	imgSrc := ""
	for _, a := range attrs {
		if a.Key == "src" {
			imgSrc = a.Val
		}
	}
	// get num param
	num := strings.SplitN(imgSrc, "num=", 2)[1]
	if num != "" {
		num := strings.SplitN(num, "&", 2)[0]
		val, err := strconv.Atoi(strings.Replace(num, ".", "", -1))
		if err == nil {
			return val
		}
	}
	return 0
}

// Returns the list of players contained in the specified page.
func ParsePlayerList(playerListPageResponse string) ([]*model.Player, []error) {
	var players []*model.Player
	var errors []error
	doc, err := gosoup.Parse(strings.NewReader(playerListPageResponse))
	if err != nil {
		return players, append(errors, err)
	}
	body := doc.DescendantsByTag("body").First()
	elts := body.DescendantsByAttrValueContaining("href", "main/fiche&voirpseudo=").All()
	for _, elt := range elts {
		usernameCell := elt.Parent
		assertIsTag(usernameCell, "td")
		userRow := usernameCell.Parent
		assertIsTag(userRow, "tr")
		player, err := parsePlayer(userRow)
		if err != nil {
			errors = append(errors, err)
		} else {
			players = append(players, player)
		}
	}
	return players, errors
}

// Creates a new player from the cells in the specified {@code <tr>} element.
func parsePlayer(playerRow *gosoup.Node) (*model.Player, error) {
	assertIsTag(playerRow, "tr")
	fields := playerRow.ChildrenByTag("td").All()
	player := new(model.Player)

	// name
	nameLink := fields[2].Children().First()
	player.Name = nameLink.Children().First().Data
	if player.Name == "" {
		player.Name = "player"
	}

	// rank
	rankElt := fields[0].Children().First()
	val, err := getAsInt(rankElt)
	if err != nil {
		return nil, fmt.Errorf("cannot parse %s's rank: %v", player.Name, err)
	}
	player.Rank = val

	// gold
	goldElt := fields[3].Children().First()
	if goldElt.Type == gosoup.TextNode {
		// gold amount is textual
		val, err := getAsInt(goldElt)
		if err != nil {
			return nil, fmt.Errorf("cannot parse %s's gold: %v", player.Name, err)
		}
		player.Gold = val
	} else {
		// gold amount is an image
		val, err := getGoldFromImgElement(goldElt)
		if err != nil {
			return nil, fmt.Errorf("OCR failed on %s's gold: %v", player.Name, err)
		}
		player.Gold = val
	}

	// army
	armyElt := fields[4]
	army := armyElt.Children().First().Data
	player.Army = model.GetArmy(army)

	// alignment
	alignmentElt := fields[5].Children().First()
	alignment := alignmentElt.Children().First().Data
	player.Alignment = model.GetAlignment(alignment)

	return player, nil
}

func getAsInt(n *gosoup.Node) (int, error) {
	numberStr := strings.Replace(n.Data, ".", "", -1)
	value, err := strconv.Atoi(numberStr)
	if err != nil {
		return value, fmt.Errorf("cannot parse %q as an int", n.Data)
	}
	return value, nil
}

// Gets the amount of stolen golden from the attack report.
func ParseGoldStolen(attackReportResponse string) (int, error) {
	doc, err := gosoup.Parse(strings.NewReader(attackReportResponse))
	if err != nil {
		return -1, err
	}
	body := doc.DescendantsByTag("body").First()
	elts := body.DescendantsByAttrValueContaining("class", "combat_gagne").All()
	if len(elts) == 0 {
		elts := body.DescendantsByAttrValueContaining("class", "combat_perdu").All()
		if len(elts) > 0 {
			return 0, nil // battle lost
		} else {
			return -1, fmt.Errorf("the page does not seem right: battle neither lost nor won")
		}
	}
	divVictory := elts[0].Parent.Parent
	return getAsInt(divVictory.DescendantsByTag("b").First().Children().First())
}

// ParseWeaponsWornness parses the weapons page response and returns the percentage
// of wornness of the weapons.
func ParseWeaponsWornness(weaponsPageResponse string) (int, error) {
	doc, err := gosoup.Parse(strings.NewReader(weaponsPageResponse))
	if err != nil {
		return -1, err
	}
	body := doc.DescendantsByTag("body").First()
	input := body.DescendantsByAttrValueContaining("title", "Armes endommagées").First()
	value := input.Children().First().Data
	if strings.HasSuffix(value, "%") {
		return strconv.Atoi(strings.TrimSuffix(value, "%"))
	} else {
		return -1, fmt.Errorf("the page does not seem right: missing %")
	}
}

// ParsePlayerGold parses the amount of gold on the specified player page.
func ParsePlayerGold(playerPageResponse string) (int, error) {
	doc, err := gosoup.Parse(strings.NewReader(playerPageResponse))
	if err != nil {
		return -1, err
	}
	body := doc.DescendantsByTag("body").First()
	img := body.DescendantsByAttrValueContaining("src", "aff_montant").First()
	return getGoldFromImgElement(img)
}

// Uses an OCR to recognize a number in the specified {@code <img>} element.
func getGoldFromImgElement(goldImageElement *gosoup.Node) (int, error) {
	assertIsTag(goldImageElement, "img")
	goldImgUrl := goldImageElement.Attr("src")
	if len(goldImgUrl) == 0 {
		panic("emtpy gold image url")
	}
	img, err := ocr.GetImage(BASE_URL + goldImgUrl)
	if err != nil {
		return -1, err
	}
	return ocr.ReadValue(&img)
}

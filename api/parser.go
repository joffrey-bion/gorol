package api

import (
	"code.google.com/p/go-charset/charset"
	_ "code.google.com/p/go-charset/data"
	//	"golang.org/x/net/html/charset"
	"errors"
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
		return "", errors.New("String: error reading the response:" + err.Error())
	}
	// get the charset from the incorrect html (hopefully not that incorrect)
	cset, err := gosoup.GetDocCharset(doc)
	if err != nil {
		logger.Println(err)
		return "", errors.New("String: error getting the charset" + err.Error())
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
		return "", errors.New("String: error reading the response with the appropriate charset")
	}
	// read the UTF-8 data as a String
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		logger.Println(err)
		return "", errors.New("String: error reading the response:" + err.Error())
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
func ParsePlayerList(playerListPageResponse string) ([]model.Player, error) {
	doc, err := gosoup.Parse(strings.NewReader(playerListPageResponse))
	if err != nil {
		return nil, err
	}
	body := doc.DescendantsByTag("body").First()
	elts := body.DescendantsByAttrValueContaining("href", "main/fiche&voirpseudo=").All()
	var list []model.Player
	for _, elt := range elts {
		usernameCell := elt.Parent
		assertIsTag(usernameCell, "td")
		userRow := usernameCell.Parent
		assertIsTag(userRow, "tr")
		list = append(list, parsePlayer(userRow))
	}
	return list, nil
}

// Creates a new player from the cells in the specified {@code <tr>} element.
func parsePlayer(playerRow *gosoup.Node) model.Player {
	assertIsTag(playerRow, "tr")
	fields := playerRow.ChildrenByTag("td").All()
	player := new(model.Player)

	// rank
	rankElt := fields[0].Children().First()
	player.Rank = getAsInt(rankElt)

	// name
	nameLink := fields[2].Children().First()
	player.Name = nameLink.Children().First().Data

	// gold
	goldElt := fields[3].Children().First()
	if goldElt.Type == gosoup.TextNode {
		// gold amount is textual
		player.Gold = getAsInt(goldElt)
	} else {
		// gold amount is an image
		player.Gold = getGoldFromImgElement(goldElt)
	}

	// army
	armyElt := fields[4]
	army := armyElt.Children().First().Data
	player.Army = model.GetArmy(army)

	// alignment
	alignmentElt := fields[5].Children().First()
	alignment := alignmentElt.Children().First().Data
	player.Alignment = model.GetAlignment(alignment)

	return *player
}

func getAsInt(n *gosoup.Node) int {
	numberStr := strings.Replace(n.Data, ".", "", -1)
	value, err := strconv.Atoi(numberStr)
	if err != nil {
		panic(fmt.Sprintf("cannot parse %q as an int", n.Data))
	}
	return value
}

//    /**
//     * Gets the amount of stolen golden from the attack report.
//     *
//     * @param attackReportResponse
//     *            the response containing the attack report.
//     * @return the amount of stolen gold, or -1 if the report couldn't be read properly
//     */
//    public static int parseGoldStolen(String attackReportResponse) {
//        Element body = Jsoup.parse(attackReportResponse).body()
//        Elements elts = body.getElementsByAttributeValue("class", "combat_gagne")
//        if (elts.size() == 0) {
//            return -1;
//        }
//        Element divVictory = elts.get(0).parent().parent()
//        return getTextAsNumber(divVictory.getElementsByTag("b").get(0))
//    }
//
//    /**
//     * Gets and parses the text contained in the spacified {@link Element}.
//     *
//     * @param numberElement
//     *            an element containing a text representing an integer, with possible dots as
//     *            thousand separator.
//     * @return the parsed number
//     */
//    private static int getTextAsNumber(Element numberElement) {
//        String number = numberElement.text().trim()
//        return Integer.valueOf(number.replace(".", ""))
//    }
//
//    /**
//     * Parses the weapons page response to return the current state of the weapons.
//     *
//     * @param weaponsPageResponse
//     *            weapons page response
//     * @return the current percentage of wornness of the weapons
//     */
//    public static int parseWeaponsWornness(String weaponsPageResponse) {
//        Element body = Jsoup.parse(weaponsPageResponse).body()
//        Elements elts = body.getElementsByAttributeValueContaining("title", "Armes endommagées")
//        Element input = elts.get(0)
//        String value = input.text().trim()
//        if (value.endsWith("%")) {
//            return Integer.valueOf(value.substring(0, value.length() - 1))
//        } else {
//            return -1
//        }
//    }
//
//    /**
//     * Parses the amount of gold on the specified player page.
//     *
//     * @param playerPageResponse
//     *            the details page of a player
//     * @return the amount of gold parsed, or -1 if it couldn't be parsed
//     */
//    public static int parsePlayerGold(String playerPageResponse) {
//        Element body = Jsoup.parse(playerPageResponse).body()
//        Elements elts = body.getElementsByAttributeValueContaining("src", "aff_montant")
//        Element img = elts.get(0)
//        return getGoldFromImgElement(img)
//    }

// Uses an OCR to recognize a number in the specified {@code <img>} element.
func getGoldFromImgElement(goldImageElement *gosoup.Node) int {
	assertIsTag(goldImageElement, "img")
	goldImgUrl := goldImageElement.Attr("src")
	if len(goldImgUrl) == 0 {
		panic("emtpy gold image url")
	}
	img, err := ocr.GetImage(BASE_URL + goldImgUrl)
	if err != nil {
		return 0
	}
	val, _ := ocr.ReadValue(&img)
	return val
}

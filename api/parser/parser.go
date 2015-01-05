package parser

import (
	"code.google.com/p/go-charset/charset"
	"errors"
	"github.com/joffrey-bion/gosoup"
	"github.com/joffrey-bion/gorol/model"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	logger *log.Logger = log.New(os.Stderr, "parser: ", 0)
)

func getCharset(resp *http.Response) string {
	if len(resp.TransferEncoding) > 0 {
		return resp.TransferEncoding[0]
	}
	doc, err := html.Parse(resp.Body)
	if err != nil {
		logger.Println("error parsing the response as HTML document")
		return ""
	}
	return gosoup.GetDocCharset(doc)
}

func readAsString(resp *http.Response) (string, error) {
	defer resp.Body.Close()

	reader, err := charset.NewReader(getCharset(resp), resp.Body)
	if err != nil {
		logger.Println(err)
		return "", errors.New("error reading the response with the appropriate charset")
	}
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		logger.Println(err)
		return "", errors.New("error reading the response")
	}

	return string(body), nil
}

func Contains(resp *http.Response, s string) bool {
	body, err := readAsString(resp)
	if err != nil {
		logger.Println(err)
		return false
	}
	logger.Println(body)
	return strings.Contains(body, s)
}

// Updates the specified account state based on the top elements of the specified page.
func UpdateState(state *model.AccountState, respReader io.Reader) {
	doc, err := html.Parse(respReader)
	if err != nil {
		return
	}
	state.Gold = findNumValueInImgUrl(doc, "onmouseover", "A chaque tour de jeu")
	state.ChestGold = findNumValueInImgUrl(doc, "onmouseover", "Votre coffre magique")
	state.Mana = findNumValueInImgUrl(doc, "onmouseover", "Votre mana représente")
	state.Turns = findNumValueInImgUrl(doc, "onmouseover", "Un nouveau tour de jeu")
	state.Adventurins = findNumValueInImgUrl(doc, "href", "main/aventurines_detail")
}

func findNumValueInImgUrl(node *html.Node, attrKey, attrValuePart string) int {
	ch := gosoup.GetChildrenByAttributeValueContaining(node, attrKey, attrValuePart)
	attrs := (<-ch).FirstChild.Attr
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

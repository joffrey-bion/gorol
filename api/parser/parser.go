package parser

import (
	"code.google.com/p/go-charset/charset"
	_ "code.google.com/p/go-charset/data"
//	"golang.org/x/net/html/charset"
	"errors"
	"github.com/joffrey-bion/gosoup"
	"github.com/joffrey-bion/gorol/model"
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

func readAsString(resp *http.Response) (string, error) {
	defer resp.Body.Close()
	doc, err := gosoup.Parse(resp.Body)
	if err != nil {
		return "", err
	}
	cset, err := gosoup.GetDocCharset(doc)
	//contentType, err := gosoup.GetDocContentType(doc)
	if err != nil {
		logger.Println(err)
		return "", errors.New("readAsString: error reading the charset")
	}
	reader, err := charset.NewReader(cset, resp.Body)
	//reader, err := charset.NewReader(resp.Body, contentType)
	if err != nil {
		logger.Println(err)
		return "", errors.New("readAsString: error reading the response with the appropriate charset")
	}
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		logger.Println(err)
		return "", errors.New("readAsString: error reading the response:" + err.Error())
	}

	return string(body), nil
}

func Contains(resp *http.Response, s string) bool {
	body, err := readAsString(resp)
	if err != nil {
		logger.Println(err)
		return false
	}
	logger.Println("Contains: " + body)
	return strings.Contains(body, s)
}

// Updates the specified account state based on the top elements of the specified page.
func UpdateState(state *model.AccountState, respReader io.Reader) {
	doc, err := gosoup.Parse(respReader)
	if err != nil {
		return
	}
	state.Gold = findNumValueInImgUrl(doc, "onmouseover", "A chaque tour de jeu")
	state.ChestGold = findNumValueInImgUrl(doc, "onmouseover", "Votre coffre magique")
	state.Mana = findNumValueInImgUrl(doc, "onmouseover", "Votre mana représente")
	state.Turns = findNumValueInImgUrl(doc, "onmouseover", "Un nouveau tour de jeu")
	state.Adventurins = findNumValueInImgUrl(doc, "href", "main/aventurines_detail")
}

func findNumValueInImgUrl(node *gosoup.Node, attrKey, attrValuePart string) int {
	ch, exit := node.ChildrenByAttrValueContaining(attrKey, attrValuePart)
	attrs := (<-ch).FirstChild.Attrs
	exit <- true
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

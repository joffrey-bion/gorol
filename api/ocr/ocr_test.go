package ocr

import (
	"fmt"
	"testing"
)

var (
	TEST_DIGITS []digitImage = []digitImage{
		{"http://s30.postimg.org/6smeyo4vx/image.png", "0"},
		{"http://s30.postimg.org/cjclcecvx/image.png", "1"},
		{"http://s30.postimg.org/pm87vo33x/image.png", "2"},
		{"http://s30.postimg.org/yjsxt0vjx/image.png", "3"},
		{"http://s30.postimg.org/5ioljmb3x/image.png", "4"},
		{"http://s30.postimg.org/6z041re0t/image.png", "5"},
		{"http://s30.postimg.org/no1jxoam5/image.png", "6"},
		{"http://s30.postimg.org/fxursj8al/image.png", "7"},
		{"http://s30.postimg.org/heavof0l9/image.png", "8"},
		{"http://s30.postimg.org/5miezvgl9/image.png", "9"},
		{"http://s30.postimg.org/5qgttvbgd/image.png", ""}} // dot
	TEST_IMGS []testImage = []testImage{
		{"http://www.riseoflords.com/scripts/aff_montant.php?montant=BTdXPAErUWxSYww7VXtSNl0%2FAWQ%3D",
			11901450,
			[]int{0, 6, 12, 14, 15, 16, 22, 28, 29, 35, 37, 38, 39, 45, 51, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69}},
		{"http://www.riseoflords.com/scripts/aff_montant.php?montant=VmZVO1JnUXsHMVJtCDs%3D",
			341793,
			[]int{0, 6, 12, 18, 19, 21, 22, 28, 29, 35, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69}},
		{"http://www.riseoflords.com/scripts/aff_montant.php?montant=Dz1UPAQ4UXtVZgYwVWU%3D",
			128200,
			[]int{0, 6, 12, 18, 19, 21, 22, 28, 29, 35, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69}},
		{"http://www.riseoflords.com/scripts/aff_montant.php?montant=UWBSOlV%2FCjdXbgM8CCZWNlw7UTQ%3D",
			22989000,
			[]int{0, 6, 12, 14, 15, 16, 22, 28, 29, 35, 37, 38, 39, 45, 51, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69}}}
)

type digitImage struct {
	URL   string
	Value string
}

type testImage struct {
	URL       string
	Value     int
	EmptyCols []int
}

func assertEquals(t *testing.T, got, expected int) {
	if got != expected {
		t.Fatalf("got %d, expected %d", got, expected)
	}
}

func TestGetEmptyColumns(t *testing.T) {
	for _, testCase := range TEST_IMGS {
		img, err := GetImage(testCase.URL)
		if err != nil {
			t.Fatalf("couldn't get an image from the http response: %v", err)
		}
		emptyCols := getEmptyColumns(&img)
		for _, col := range testCase.EmptyCols {
			if !emptyCols[col] {
				t.Fatalf("column %d is supposed to be empty!", col)
			}
		}
	mainLoop:
		for col, empty := range emptyCols {
			if empty {
				for _, tcol := range testCase.EmptyCols {
					if col == tcol {
						// we found the column among the true empty columns
						continue mainLoop
					}
				}
				// we didn't find the column among the true empty columns
				t.Fatalf("column %d is NOT supposed to be empty!", col)
			}
		}
	}
}

func TestSimilarDigits(t *testing.T) {
	for _, testCase := range TEST_DIGITS {
		img, err := GetImage(testCase.URL)
		if err != nil {
			t.Fatalf("couldn't get an image from the http response: %v", err)
		}
		str, err := getDigit(&img)
		if err != nil {
			fmt.Printf("Image preview for %d:\n", testCase.Value)
			printAsciiImage(&img)
			t.Fatalf("couldn't recognized the digit %q", testCase.Value)
		}
		if str != testCase.Value {
			fmt.Printf("Image preview for %d:\n", testCase.Value)
			printAsciiImage(&img)
			t.Fatalf("digit %q expected, got %q instead", testCase.Value, str)
		}
	}
}

func TestOCR(t *testing.T) {
	for _, testCase := range TEST_IMGS {
		img, err := GetImage(testCase.URL)
		if err != nil {
			t.Fatalf("couldn't get an image from the http response: %v", err)
		}
		val, err := ReadValue(&img)
		if err != nil {
			fmt.Printf("Image preview for %d:\n", testCase.Value)
			printAsciiImage(&img)
			t.Fatalf("couldn't read the value (expected %d) of the image: %v", testCase.Value, err)
		}
		assertEquals(t, val, testCase.Value)
	}
}

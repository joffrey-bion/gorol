package ocr

import (
    "testing"
)

func assertEquals(t *testing.T, got, expected int) {
	if got != expected {
		t.Fatalf("got %d, expected %d", got, expected)
	}
}

func assertImage(t *testing.T, expectedValue int, url string) {
	img, err := GetImage("http://www.riseoflords.com/scripts/aff_montant.php?montant=Dj8AaVF7AzcMO1ZmCCZdOQdgA2Y%3D")
	if err != nil {
		t.Fatalf("couldn't get an image from the http response: %v", err)
	}
	val, err := ReadValue(img)
	if err != nil {
		t.Fatalf("couldn't read the value of the image: %v", err)
	}
	assertEquals(t, val, expectedValue)
}

func TestOCR(t *testing.T) {
	assertImage(t, 11901450, "http://www.riseoflords.com/scripts/aff_montant.php?montant=BTdXPAErUWxSYww7VXtSNl0%2FAWQ%3D")
	
	assertImage(t, 341793, "http://www.riseoflords.com/scripts/aff_montant.php?montant=VmZVO1JnUXsHMVJtCDs%3D")
	
	assertImage(t, 128200, "http://www.riseoflords.com/scripts/aff_montant.php?montant=Dz1UPAQ4UXtVZgYwVWU%3D")
	
	assertImage(t, 22989000, "http://www.riseoflords.com/scripts/aff_montant.php?montant=UWBSOlV%2FCjdXbgM8CCZWNlw7UTQ%3D")
}


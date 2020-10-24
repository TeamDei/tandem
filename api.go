package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

const (
	API_Base = "http://www.perseus.tufts.edu/hopper/xmlmorph?lang=la&lookup="
)

// Response stores analyses on a Latin word
type Response struct {
	Analyses []Analysis `xml:"analysis"`
}

// Analysis stores an interpretation of a Latin word
type Analysis struct {
	Form string `xml:"form"`
	Lemma string `xml:"lemma"`
	ExpandedForm string `xml:"expandedForm"`
	Pos string `xml:"pos"`
	Number string `xml:"number"`
	Gender string `xml:"gender"`
	Case string `xml:"case"`
	Dialect string `xml:"dialect"`
	Feature string `xml:"feature"`
	Person string `xml:"person"`
	Tense string `xml:"tense"`
	Mood string `xml:"mood"`
	Voice string `xml:"voice"`
}

// String returns a pretty-printed representation of an Analysis
func (a *Analysis) String() string {
	switch a.Pos {
	case "noun":
		fallthrough
	case "pron":
		fallthrough
	case "adj":
		if a.Feature != "" {
			a.Feature = " ("+a.Feature+")"
		}
		return fmt.Sprintf("(%s)%s %s. %s. %s. of %s", a.Pos, a.Feature, a.Case, a.Gender, a.Number, a.Lemma)
	case "verb":
		if a.Feature != "" {
			a.Feature = "("+a.Feature+")"
		}
		if a.Person != "" {
			a.Person = " "+a.Person+" person "
			a.Number += ","
		}
		return fmt.Sprintf("(%s)%s%s%s %s. %s. %s. of %s", a.Pos, a.Feature, a.Person, a.Number, a.Voice, a.Mood, a.Tense, a.Lemma)
	case "adv":
		return fmt.Sprintf("(%v) %v", a.Pos, a.Lemma)
	default:
		return fmt.Sprintf("error, cannot convert %s to string", a.Pos)
	}
}

// Generate response from the API
func GenAPI(word string) (Response, error) {
	r, err := http.Get(API_Base+word)
	if err != nil {
		return Response{}, err
	}
	defer r.Body.Close()
	return generateResponse(r.Body)
}

// Generate response from a file
func GenFile(path string) (Response, error) {
	f, err := os.Open(path)
	if err != nil {
		return Response{}, err
	}
	defer f.Close()
	return generateResponse(f)
}

// Generate response
func generateResponse(f io.Reader) (Response, error) {
	body, err := ioutil.ReadAll(f)
	if err != nil {
		return Response{}, err
	}
	var resp Response
	if err = xml.Unmarshal(body, &resp); err != nil {
		return Response{}, err
	} else {
		return resp, nil
	}
}

// Perform tests (optionally test the API too)
func Tests(API bool) {
	Test("noun");Test("pronoun");Test("verb");Test("adjective");Test("adverb")
	if API {
		TestAPI("nautam");TestAPI("noster");TestAPI("hortari");TestAPI("acri");TestAPI("diligenter")
	}
	fmt.Println("")
}

// Parse an example xml and print analyses
func Test(id string) {
	fmt.Printf("\n===%s test===\n", id)
	r, err := GenFile(path.Join("test_data", "example_"+id+".xml"))
	if err != nil {
		panic(err)
	}
	for _, a := range r.Analyses {
		fmt.Println(a.String())
	}
}

// Parse example words from the actual API and print analyses
func TestAPI(word string) {
	fmt.Printf("\n===API TEST (%s)===\n", word)
	r, err := GenAPI(word)
	if err != nil {
		panic(err)
	}
	for _, a := range r.Analyses {
		fmt.Println(a.String())
	}
}

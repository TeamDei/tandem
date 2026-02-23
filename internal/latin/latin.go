package latin

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
)

const APIBase = "http://www.perseus.tufts.edu/hopper/xmlmorph?lang=la&lookup="

// punctRe matches any character that should be stripped before an API lookup:
// carriage returns, newlines, and all POSIX punctuation (covers .,;!@#$%^&*()-_+
// and more).
var punctRe = regexp.MustCompile(`[\r\n[:punct:]]`)

// Response stores analyses on a Latin word.
type Response struct {
	Analyses []Analysis `xml:"analysis"`
}

// Analysis stores one morphological interpretation of a Latin word.
type Analysis struct {
	Form         string `xml:"form"`
	Lemma        string `xml:"lemma"`
	ExpandedForm string `xml:"expandedForm"`
	Pos          string `xml:"pos"`
	Number       string `xml:"number"`
	Gender       string `xml:"gender"`
	Case         string `xml:"case"`
	Dialect      string `xml:"dialect"`
	Feature      string `xml:"feature"`
	Person       string `xml:"person"`
	Tense        string `xml:"tense"`
	Mood         string `xml:"mood"`
	Voice        string `xml:"voice"`
}

// String returns a human-readable representation of the analysis.
func (a Analysis) String() string {
	switch a.Pos {
	case "noun", "pron", "conj", "prep", "adj":
		feature := ""
		if a.Feature != "" {
			feature = " (" + a.Feature + ")"
		}
		return fmt.Sprintf("(%s)%s %s. %s. %s. of %s", a.Pos, feature, a.Case, a.Gender, a.Number, a.Lemma)
	case "part", "verb":
		feature := ""
		if a.Feature != "" {
			feature = "(" + a.Feature + ")"
		}
		person := ""
		number := a.Number
		if a.Person != "" {
			person = " " + a.Person + " person "
			number += ","
		}
		return fmt.Sprintf("(%s)%s%s%s %s. %s. %s. of %s", a.Pos, feature, person, number, a.Voice, a.Mood, a.Tense, a.Lemma)
	case "adv":
		return fmt.Sprintf("(%v) %v", a.Pos, a.Lemma)
	default:
		return fmt.Sprintf("error, cannot convert %s to string", a.Pos)
	}
}

// Lookup queries the Perseus API for morphological analyses of word.
func Lookup(word string) (Response, error) {
	word = punctRe.ReplaceAllString(word, "")
	r, err := http.Get(APIBase + word)
	if err != nil {
		return Response{}, err
	}
	defer func() {
		_ = r.Body.Close()
	}()
	return parseResponse(r.Body)
}

// LookupFile parses morphological analyses from a local XML file.
func LookupFile(filePath string) (Response, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return Response{}, err
	}
	defer func() {
		_ = f.Close()
	}()
	return parseResponse(f)
}

// parseResponse decodes an XML morphology response from r.
func parseResponse(r io.Reader) (Response, error) {
	body, err := io.ReadAll(r)
	if err != nil {
		return Response{}, err
	}
	var resp Response
	if err = xml.Unmarshal(body, &resp); err != nil {
		return Response{}, err
	}
	return resp, nil
}

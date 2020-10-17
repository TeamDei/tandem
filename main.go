package main

import (
	"fmt"
	"path"
)

func main() {
	Test("noun");Test("pronoun");Test("verb");Test("adjective");Test("adverb")
	//TestAPI("nautam");TestAPI("noster");TestAPI("hortari");TestAPI("acri");TestAPI("diligenter")
	fmt.Println("")
}

// Parse an example and print analyses
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

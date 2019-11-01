package main

import (
	"fmt"

	"github.com/athoune/bleve-lexicon/lexicon"
	"github.com/blevesearch/bleve"
)

func main() {
	// open a new index
	mapping := bleve.NewIndexMapping()
	index, err := bleve.New("example.bleve", mapping)
	if err != nil {
		fmt.Println(err)
		return
	}

	// index some data
	index.Index("1", struct {
		Name string
	}{
		Name: "Il ne faut pas prendre les enfants du bon Dieu pour des canards sauvages",
	})
	index.Index("2", struct {
		Name string
	}{
		Name: "la cité des enfants perdus",
	})
	index.Index("3", struct {
		Name string
	}{
		Name: "Les enfants du paradis",
	})
	index.Index("4", struct {
		Name string
	}{
		Name: "Rats des villes et rats des champs",
	})

	l, err := lexicon.Lexicon(index, "_all")
	if err != nil {
		panic(err)
	}
	f, err := l.FieldDict("_all")
	if err != nil {
		panic(err)
	}
	for {
		e, err := f.Next()
		if err != nil || e == nil {
			break
		}
		fmt.Println(e.Term, e.Count)
	}
}

package main

import (
	"fmt"

	"github.com/athoune/bleve-lexicon/lexicon"
	"github.com/blevesearch/bleve"
)

func main() {
	// open a new index
	mapping := bleve.NewIndexMapping()
	index, err := lexicon.OpenOrCreate("example.bleve", mapping)
	if err != nil {
		panic(err)
	}

	batch := index.NewBatch()
	// index some data
	batch.Index("1", struct {
		Name string
	}{
		Name: "Il ne faut pas prendre les enfants du bon Dieu pour des canards sauvages",
	})
	batch.Index("2", struct {
		Name string
	}{
		Name: "la cité des enfants perdus",
	})
	batch.Index("3", struct {
		Name string
	}{
		Name: "Les enfants du paradis",
	})
	batch.Index("4", struct {
		Name string
	}{
		Name: "Rats des villes et rats des champs",
	})

	err = index.Batch(batch)
	if err != nil {
		panic(err)
	}
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

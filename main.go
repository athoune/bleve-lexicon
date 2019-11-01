package main

import (
	"fmt"

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

	data := struct {
		Name string
	}{
		Name: "Il ne faut pas prendre les enfants du bon Dieu pour des canards sauvages",
	}

	// index some data
	index.Index("id", data)

	fmt.Println(index.StatsMap())
	fmt.Println(index.Fields())

	dict, err := index.FieldDict("_all")
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		entry, err := dict.Next()
		if err != nil || entry == nil {
			break
		}
		fmt.Println(entry)
	}

	_, kv, err := index.Advanced()
	if err != nil {
		fmt.Println(err)
		return
	}
	reader, err := kv.Reader()
	if err != nil {
		fmt.Println(err)
		return
	}

	iter := reader.PrefixIterator([]byte(""))
	for {
		k, v, ok := iter.Current()
		fmt.Println(string(k), "\n\t", string(v))
		if !ok {
			break
		}
		iter.Next()
	}

	// search for some text
	query := bleve.NewMatchQuery("text")
	search := bleve.NewSearchRequest(query)
	searchResults, err := index.Search(search)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(searchResults)
}

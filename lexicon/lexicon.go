package lexicon

import (
	"math"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/custom"
	"github.com/blevesearch/bleve/analysis/char/asciifolding"
	"github.com/blevesearch/bleve/analysis/token/ngram"
	"github.com/blevesearch/bleve/analysis/tokenizer/single"
	"github.com/blevesearch/bleve/mapping"
	"github.com/blevesearch/bleve/search"
)

type token struct {
	Value string
}

var LexiconMapping mapping.IndexMapping

func init() {
	mapping := bleve.NewIndexMapping()
	err := mapping.AddCustomTokenFilter(ngram.Name, map[string]interface{}{
		"min":  3,
		"max":  3,
		"type": ngram.Name,
	})
	if err != nil {
		panic(err)
	}
	err = mapping.AddCustomCharFilter(asciifolding.Name, map[string]interface{}{
		"type": asciifolding.Name,
	})
	if err != nil {
		panic(err)
	}
	err = mapping.AddCustomAnalyzer("custom1", map[string]interface{}{
		"type":          custom.Name,
		"tokenizer":     single.Name,
		"token_filters": []interface{}{ngram.Name},
		"char_filters":  []interface{}{asciifolding.Name},
	})
	if err != nil {
		panic(err)
	}
	mapping.DefaultAnalyzer = "custom1"
	LexiconMapping = mapping

}

type Lexicon struct {
	Index bleve.Index
}

// Lexicon returns an index of tokens
func New(index bleve.Index, name string) (*Lexicon, error) {
	lexicon, err := OpenOrCreate("lexicon.bleve", LexiconMapping)
	if err != nil {
		return nil, err
	}
	batch := lexicon.NewBatch()
	dict, err := index.FieldDict(name)
	if err != nil {
		return nil, err
	}
	for {
		entry, err := dict.Next()
		if err != nil || entry == nil {
			break
		}
		err = batch.Index(entry.Term, token{Value: entry.Term})
		if err != nil {
			return nil, err
		}
	}
	err = lexicon.Batch(batch)
	if err != nil {
		return nil, err
	}
	return &Lexicon{
		Index: lexicon,
	}, nil
}

// DoYouMean suggest an other spelling
func (l *Lexicon) DoYouMean(word string, maxDiff int) ([]string, error) {
	query := bleve.NewMatchQuery(word)
	query.SetField("Value")
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Fields = []string{"*"}
	searchRequest.Explain = true
	searchResult, err := l.Index.Search(searchRequest)
	if err != nil {
		return nil, err
	}
	r := make([]string, 0)
	for _, h := range searchResult.Hits {
		if math.Abs(float64(len(h.ID)-len(word))) > float64(maxDiff) {
			continue
		}
		if search.LevenshteinDistance(h.ID, word) > maxDiff {
			continue
		}
		r = append(r, h.ID)
	}
	return r, nil
}

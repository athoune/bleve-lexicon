package lexicon

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/custom"
	"github.com/blevesearch/bleve/analysis/char/asciifolding"
	"github.com/blevesearch/bleve/analysis/token/ngram"
	"github.com/blevesearch/bleve/analysis/tokenizer/single"
)

type token struct {
	Value string
}

// Lexicon returns an index of tokens
func Lexicon(index bleve.Index, name string) (bleve.Index, error) {
	mapping := bleve.NewIndexMapping()
	err := mapping.AddCustomTokenFilter(ngram.Name, map[string]interface{}{
		"min":  3,
		"max":  3,
		"type": ngram.Name,
	})
	if err != nil {
		panic(err)
		return nil, err
	}
	err = mapping.AddCustomCharFilter(asciifolding.Name, map[string]interface{}{
		"type": asciifolding.Name,
	})
	if err != nil {
		panic(err)
		return nil, err
	}
	err = mapping.AddCustomAnalyzer("custom1", map[string]interface{}{
		"type":          custom.Name,
		"tokenizer":     single.Name,
		"token_filters": []interface{}{ngram.Name},
		"char_filters":  []interface{}{asciifolding.Name},
	})
	if err != nil {
		panic(err)
		return nil, err
	}
	mapping.DefaultAnalyzer = "custom1"
	lexicon, err := bleve.New("lexicon.bleve", mapping)
	if err != nil {
		panic(err)
		return nil, err
	}
	dict, err := index.FieldDict(name)
	if err != nil {
		return nil, err
	}
	for {
		entry, err := dict.Next()
		if err != nil || entry == nil {
			break
		}
		lexicon.Index(entry.Term, token{Value: entry.Term})
	}
	return lexicon, nil
}

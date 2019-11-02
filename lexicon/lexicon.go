package lexicon

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/custom"
	"github.com/blevesearch/bleve/analysis/char/asciifolding"
	"github.com/blevesearch/bleve/analysis/token/ngram"
	"github.com/blevesearch/bleve/analysis/tokenizer/single"
	"github.com/blevesearch/bleve/mapping"
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

// Lexicon returns an index of tokens
func Lexicon(index bleve.Index, name string) (bleve.Index, error) {
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
	return lexicon, nil
}

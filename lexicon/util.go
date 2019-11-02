package lexicon

import (
	"os"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/index/scorch"
	"github.com/blevesearch/bleve/mapping"
)

// OpenOrCreate try to open an index, or create it
func OpenOrCreate(path string, mapping mapping.IndexMapping) (bleve.Index, error) {
	_, err := os.Stat(path)
	if err == nil {
		return bleve.Open(path)
	}
	if os.IsNotExist(err) {
		return bleve.NewUsing(path, mapping, scorch.Name, scorch.Name, nil)
	}
	return nil, err
}

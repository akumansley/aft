package db

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
)

type Index interface {
	Init()
	Index(Ider)
	Query(string) *bleve.SearchResult // TODO wrap SearchResult
}

type BleveIndex struct {
	mapping mapping.IndexMapping
	index   bleve.Index
}

func (i *BleveIndex) Init() {
	i.mapping = bleve.NewIndexMapping()
	i.index, _ = bleve.NewMemOnly(i.mapping)
}

func (i *BleveIndex) Index(item Ider) {
	i.index.Index(item.GetId(), item)
}

func (i *BleveIndex) Query(q string) *bleve.SearchResult {
	query := bleve.NewMatchQuery(q)
	searchRequest := bleve.NewSearchRequest(query)
	searchResults, _ := i.index.Search(searchRequest)
	return searchResults
}

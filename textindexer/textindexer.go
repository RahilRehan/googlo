package textindexer

import (
	"time"

	"github.com/google/uuid"
)

type Document struct {
	LinkID    uuid.UUID
	URL       string
	Title     string
	Content   string
	IndexedAt time.Time
	PageRank  float64
}

type QueryType uint8

const (
	QueryTypeMatch  QueryType = iota // match list of keywords in any order
	QueryTypePhrase                  // exact phrase match
)

type Query struct {
	Type       QueryType
	Expression string
	Offset     uint64
}

type Indexer interface {
	Index(*Document) error
	FindByID(uuid.UUID) (*Document, error)
	Search(Query) (Iterator, error)
	UpdateScore(uuid.UUID, float64) error
}

type Iterator interface {
	Close() error // close the iterator and release all the resources
	Next() bool   //loads the next document, returns false if next doc doesn't exists
	Error() error
	Document() *Document
	TotalCount() uint64 // returns the number of search results queried
}

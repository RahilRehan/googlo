# TextIndexer

## Document
```golang
type Document struct{
    LinkId uuid.UUID // represents associated link
    URL string // represents associated links url
    Title string // title is extracted from html attribute `title` of particular link 
    Content string // entire content inside link
    IndexedAt time.Time // tells when Document was indexed/re-indexed
    PageRank float64 // will be populated by pagerank component for efficient indexing
}
```
- Some fields like `Title` can also help in search engine optimisation.
- Output results are sorted based on `PageRank` to display queries efficiently.

## Operations supported on Document
- `Index` or Reindex document
- Find document by `ID`
- Return iterable results for a search query, pagination shoould be supported
- Update `PageRank`

```golang
type Indexer interface {
    Index(*Document) error
    FindByID(uuid.UUID) (*Document, error)
    Search(Query) (Iterator, error)
    UpdateScore(uuid.UUID, float64) error
}
```

## Query type
```golang
type Query struct {
    Type       QueryType
    Expression string
    Offset     uint64
}
```
- the query type can be used as a composite to store details about query. 
- if you want to add new featured down the line, it becomes much easier, and search method in Indexer doesn't change doesn't change

## Elastic search
- can easily scale horizontally, hence highly available
- provides a nice `go` client `go-elastic` library to talk with elastic search REST api
- if the search instances scale horizontally, its is really helpful for read intensive operations.
- For right operations across multiple clusters, elastic search manages to keep consistency across nodes.



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


package cockroach

import (
	"github.com/RahilRehan/googlo/linkgraph/cockroachdb/sqlc"
	"github.com/RahilRehan/googlo/linkgraph/graph"
)

type linkIterator struct{
	currentIdx int
	links []sqlc.Links
	lastErr error
	latchedLink *graph.Link
}

type edgeIterator struct{
	currentIdx int
	edges []sqlc.Edges
	lastErr error
	latchedEdge *graph.Edge
}

func (i *linkIterator) Link() *graph.Link {
	return i.latchedLink
}

func (i *linkIterator) Next() bool {
	if i.lastErr != nil || i.currentIdx >= len(i.links) {
		return false
	}

	newLink := new(graph.Link)
	retrievedLink := i.links[i.currentIdx]

	newLink.ID = retrievedLink.ID
	newLink.URL = retrievedLink.Url.String
	newLink.RetrievedAt = retrievedLink.RetrievedAt.Time.UTC()

	i.latchedLink = newLink
	i.currentIdx++

	return true
}

func (i *linkIterator) Error() error{
	return i.lastErr
}

func (i *linkIterator) Close() error{
	i.currentIdx = len(i.links)
	return nil
}



func (i *edgeIterator) Edge() *graph.Edge {
	return i.latchedEdge
}

func (i *edgeIterator) Next() bool {
	if i.lastErr != nil || i.currentIdx >= len(i.edges) {
		return false
	}

	newEdge := new(graph.Edge)
	retrievedEdge := i.edges[i.currentIdx]

	newEdge.ID = retrievedEdge.ID
	newEdge.Source = retrievedEdge.Src
	newEdge.Destination = retrievedEdge.Dst
	newEdge.UpdatedAt = retrievedEdge.UpdatedAt.Time.UTC()

	i.latchedEdge = newEdge
	i.currentIdx++

	return true
}

func (i *edgeIterator) Error() error{
	return i.lastErr
}

func (i *edgeIterator) Close() error{
	i.currentIdx = len(i.edges)
	return nil
}




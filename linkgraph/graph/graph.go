package graph

import (
	"time"

	"github.com/google/uuid"
)

type Link struct{
	ID uuid.UUID
	URL string
	RetrievedAt time.Time
}

type Edge struct{
	ID uuid.UUID
	Source uuid.UUID
	Destination uuid.UUID
	UpdatedAt time.Time
}

type Iterator interface{
	// returns true if next instance is present in iterator, returns false otherwise or if an error occurs
	Next() bool

	// ex: when connection with db fails
	// with this pattern, we don't need to check error for each iteration, rather we can check once at the end
	Error() error 

	// Releases all the resources associated with the iterator
	Close() error
}

type LinkIterator interface{
	Iterator

	// returns the next link in the iterator
	Link() *Link
}

type EdgeIterator interface{
	Iterator

	// returns the next edge in the iterator
	Edge() *Edge
}

type Graph interface{
	// Upsert means to insert or update
	UpsertLink(link *Link) error 
	UpsertEdge(edge *Edge) error

	FindLink(id uuid.UUID) (*Link, error)
	
	// links within the ID range [fromId, toId)
	Links(fromId, toID uuid.UUID, retrievedBefore time.Time)
	// Return edges for which the origin link's ID is within the [fromId, toId)
	Edges(fromId, toId uuid.UUID, updatedBefore time.Time)
	
	RemoveStaleEdges(fromId uuid.UUID, updatedBefore time.Time) error
}


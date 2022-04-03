package memory

import (
	"sync"

	"github.com/RahilRehan/googlo/linkgraph/graph"

	"github.com/google/uuid"
)

// in-memory store for graph

type edgeList []uuid.UUID

type inMemoryGraph struct {
	// Allows multiple readers but just one writer to hold lock
	// This helps for safe concurrent access
	// Provides much better throughput compared to sync.Mutex for read heavy operations
	mu sync.RWMutex

	// Indexing based of uuid
	links map[uuid.UUID]*graph.Link
	edges map[uuid.UUID]*graph.Edge

	// Indexing based of url, as url is also unique
	linkURLIdx map[string]*graph.Link

	// Used to remove stale edges, contains list of edges that originate from particular link
	linkEdgeMap map[uuid.UUID]*edgeList
}

func NewInMemoryGraph() *inMemoryGraph{
	return &inMemoryGraph{
		links: make(map[uuid.UUID]*graph.Link),
		edges: make(map[uuid.UUID]*graph.Edge),
		linkURLIdx: make(map[string]*graph.Link),
		linkEdgeMap: make(map[uuid.UUID]*edgeList),
	}
}

func (g *inMemoryGraph) UpsertLink(link *graph.Link) error {
	// This is always a write operation, hence will acquire a write lock
	g.mu.Lock()
	defer g.mu.Unlock()

	// If link with same URL already exists, then update
	if existing := g.linkURLIdx[link.URL]; existing != nil {
		link.ID = existing.ID // Don't create a new ID, assign existing one
		origTs := existing.RetrievedAt 
		*existing = *link // update link
		if origTs.After(link.RetrievedAt){
			existing.RetrievedAt = origTs  // update RetrivedAt with latest timestamp
		}
		return nil
	}

	// create new id until a unique ID is created
	for {
		link.ID = uuid.New()
		if g.links[link.ID] == nil {
			break
		}
	}

	// If link doesn't already exists, create a new one
	copy := new(graph.Link)  // we create this copy to make sure, no code outside of implementation can modify graph data
	*copy = *link
	g.links[copy.ID] = copy
	g.linkURLIdx[copy.URL] = copy

	return nil
}

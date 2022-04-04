package memory

import (
	"fmt"
	"sync"
	"time"

	"github.com/RahilRehan/googlo/linkgraph/graph"

	"github.com/google/uuid"
)

// in-memory store for graph

// Compile-time check for ensuring InMemoryGraph implements Graph.
var _ graph.Graph = (*inMemoryGraph)(nil)

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

	// contains list of edges that originate from particular link, can be used to remove stale edges, 
	linkEdgeMap map[uuid.UUID]edgeList
}

func NewInMemoryGraph() *inMemoryGraph{
	return &inMemoryGraph{
		links: make(map[uuid.UUID]*graph.Link),
		edges: make(map[uuid.UUID]*graph.Edge),
		linkURLIdx: make(map[string]*graph.Link),
		linkEdgeMap: make(map[uuid.UUID]edgeList),
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

func (g *inMemoryGraph) UpsertEdge(edge *graph.Edge) error{
	// This is always a write operation, hence will acquire a write lock
	g.mu.Lock()
	defer g.mu.Unlock()

	// check if source and dest links are present
	_, srcExists := g.links[edge.Source]
	_, destExists := g.links[edge.Destination]

	if !srcExists || !destExists{
		return fmt.Errorf("upsert edge: %v", graph.ErrUnknownEdgeLinks)
	}

	// retrieve set of edges that originate from given source
	// this means that the edge with given source and destination already exists, so just update time
	for _, edgeId := range g.linkEdgeMap[edge.Source]{
		existingEdge := g.edges[edgeId]
		if existingEdge.Source == edge.Source && existingEdge.Destination == edge.Destination{
			existingEdge.UpdatedAt = time.Now()
			*edge = *existingEdge // update upstream instance to reflect changes to callers
			return nil
		}
	}

	for {
		edge.ID = uuid.New()
		if g.edges[edge.ID] == nil{
			break
		}
	}

	// create a new edge with given links
	edge.UpdatedAt = time.Now()
	copy := new(graph.Edge)
	*copy = *edge
	g.edges[copy.ID] = copy
	g.linkEdgeMap[copy.Source] = append(g.linkEdgeMap[copy.Source], copy.ID)
	
	return nil
}

func (g *inMemoryGraph) FindLink(id uuid.UUID) (*graph.Link, error){
	g.mu.RLock()
	defer g.mu.RUnlock()

	if existing := g.links[id]; existing != nil{
		copy := new(graph.Link)
		*copy = *existing
		return copy, nil
	}

	return nil, fmt.Errorf("find link: %w", graph.ErrNotFound)
}

func (g *inMemoryGraph) Links(fromId, toId uuid.UUID, retrievedBefore time.Time) (graph.LinkIterator, error){
	from, to := fromId.String(), toId.String()

	g.mu.RLock()
	defer g.mu.RUnlock()

	var list []*graph.Link
	for linkId, link := range g.links{
		if id := linkId.String(); id >= from && id < to && link.RetrievedAt.Before(retrievedBefore){
			list = append(list, link)
		}
	}

	return &linkiterator{graph: g, links: list}, nil
}

func (g *inMemoryGraph) Edges(fromId, toId uuid.UUID, updatedBefore time.Time) (graph.EdgeIterator, error){
	from, to := fromId.String(), toId.String()

	g.mu.RLock()
	defer g.mu.RUnlock()

	var list []*graph.Edge
	for edgeId, edge := range g.edges{
		if id := edgeId.String(); id >= from && id < to && edge.UpdatedAt.Before(updatedBefore){
			list = append(list, edge)
		}
	}
	return &edgeIterator{graph: g, edges: list}, nil
}

func (g *inMemoryGraph) RemoveStaleEdges(fromId uuid.UUID, updatedBefore time.Time) error{
	g.mu.Lock()
	defer g.mu.Unlock()

	var newEdgeList edgeList
	for _, edgeId := range g.linkEdgeMap[fromId]{
		edge := g.edges[edgeId]
		if edge.UpdatedAt.Before(updatedBefore){
			delete(g.edges, edgeId)
			continue
		}
		newEdgeList = append(newEdgeList, edgeId)
	}
	g.linkEdgeMap[fromId] = newEdgeList
	return nil
}

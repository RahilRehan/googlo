package cockroach

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/RahilRehan/googlo/linkgraph/cockroachdb/sqlc"
	"github.com/RahilRehan/googlo/linkgraph/graph"
)

const driver = "postgres"

type cockroachDBGraph struct{
	db *sqlc.Queries
}


func NewCockroachDBGraph(dsn string) *cockroachDBGraph {
	db, err := sql.Open(driver, dsn)
	if err != nil{
		log.Fatalln(err)
	}

	return &cockroachDBGraph{db: sqlc.New(db)}
}

func (c *cockroachDBGraph) UpsertLink(link *graph.Link) error {
	upsertedLink, err := c.db.UpsertLink(
		context.Background(), 
		sqlc.UpsertLinkParams{
			Url: sql.NullString{String: link.URL, Valid: true},
			RetrievedAt: sql.NullTime{Time: link.RetrievedAt.UTC(), Valid: true},
		},
	)

	
	if err != nil{
		return errors.New(fmt.Sprintf("upsert link: %s", err))
	}

	link.RetrievedAt = upsertedLink.RetrievedAt.Time.UTC()
	return nil
}

func (c *cockroachDBGraph) UpsertEdge(edge *graph.Edge) error{
	upsertedEdge, err := c.db.UpsertEdge(
		context.Background(),
		sqlc.UpsertEdgeParams{
			Src: edge.Source,
			Dst: edge.Destination,
		},
	)

	if isForeignKeyViolationError(err){
		return errors.New(fmt.Sprintf("upsert edge: foreign key error"))
	}

	if err != nil{
		return errors.New(fmt.Sprintf("upsert edge: %s", err))
	}

	edge.UpdatedAt = upsertedEdge.UpdatedAt.Time.UTC()
	
	return nil
}

func (c *cockroachDBGraph) FindLink(id uuid.UUID) (*graph.Link, error){
	foundLink, err := c.db.FindLink(context.Background(),id)
	if err != nil{
		return nil, errors.New(fmt.Sprintf("find link: %s", err))
	}

	return &graph.Link{
		ID: id,
		URL: foundLink.Url.String,
		RetrievedAt: foundLink.RetrievedAt.Time.UTC(),
	}, nil
}

func (c *cockroachDBGraph) Links(fromId, toId uuid.UUID, retrievedBefore time.Time) (graph.LinkIterator, error){
	links, err := c.db.LinksInPartition(
		context.Background(),
		sqlc.LinksInPartitionParams{
			ID: fromId, ID_2: toId, 
			RetrievedAt: sql.NullTime{Time: retrievedBefore, Valid:true}},
		)
	
	if err != nil{
		return nil, errors.New(fmt.Sprintf("links: %s", err))
	}

	return &linkIterator{
		currentIdx: 0,
		links: links,
	}, nil

}

func (c *cockroachDBGraph) Edges(fromId, toId uuid.UUID, updatedBefore time.Time) (graph.EdgeIterator, error){
	edges, err := c.db.EdgesInPartition(
		context.Background(),
		sqlc.EdgesInPartitionParams{
			Src: fromId,
			Dst: toId,
			UpdatedAt: sql.NullTime{Time:updatedBefore, Valid: true},
		},
	)

	if err != nil{
		return nil, errors.New(fmt.Sprintf("edges: %s", err))
	}

	return &edgeIterator{
		currentIdx: 0,
		edges: edges,
	}, nil

}

func (c *cockroachDBGraph)RemoveStaleEdges(fromId uuid.UUID, updatedBefore time.Time) error{
	err := c.db.RemoveStaleEdges(
		context.Background(),
		sqlc.RemoveStaleEdgesParams{Src: fromId, UpdatedAt: sql.NullTime{Time: updatedBefore, Valid: true}},
	)
	if err != nil{
		return errors.New(fmt.Sprintf("remove stale edges: %s", err))
	}
	return nil
}

func isForeignKeyViolationError(err error) bool {
    pqErr, valid := err.(*pq.Error)
    if !valid {
        return false
    }
    return pqErr.Code.Name() == "foreign_key_violation"
}


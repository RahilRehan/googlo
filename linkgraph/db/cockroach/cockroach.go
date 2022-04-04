package cockroach

import (
	"errors"
	"fmt"

	"github.com/RahilRehan/googlo/linkgraph/graph"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type CockroachDBGraph struct{
	db *gorm.DB
}

func NewCockroachDBGraph(dsn string) *CockroachDBGraph{
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil{
		panic(err)
	}

	return &CockroachDBGraph{db: gormDB}
}

func (g CockroachDBGraph) UpsertLink(link *graph.Link) error{
	result := g.db.Model(&graph.Link{}).Updates(link)

	if result.RowsAffected == 0{
		return fmt.Errorf("upsert link: %w", errors.New("no change"))
	}else if (result.Error != nil){
		return fmt.Errorf("upsert link: %w", result.Error)
	}

	link.RetrievedAt = link.RetrievedAt.UTC()
	return nil
}


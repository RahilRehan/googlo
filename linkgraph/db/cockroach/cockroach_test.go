package cockroach_test

import (
	"testing"

	"github.com/RahilRehan/googlo/linkgraph/db/cockroach"
	"github.com/RahilRehan/googlo/linkgraph/graph"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type cockroachDBTestSuite struct{
	suite.Suite
	cdb cockroach.CockroachDBGraph
}

func (s *cockroachDBTestSuite) SetupTest(){
	s.cdb = *cockroach.NewCockroachDBGraph("cockroachdb://root@localhost:26257/linkgraph?sslmode=disable")
}

func (s *cockroachDBTestSuite) TestUpsertLink(){
	err := s.cdb.UpsertLink(&graph.Link{
		URL: "google.com",
	})
	assert.Nil(s.T(), err)
}

func TestCockroachDB(t *testing.T){
	suite.Run(t, new(cockroachDBTestSuite))
}
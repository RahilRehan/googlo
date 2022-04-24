package cockroach_test

import (
	"testing"

	cockroach "github.com/RahilRehan/googlo/linkgraph/cockroachdb"
	"github.com/RahilRehan/googlo/linkgraph/graph"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var cockroachDBTest cockroachDBTestSuite

type cockroachDBTestSuite struct {
	suite.Suite
	cdb *cockroach.CockroachDBGraph
}

func TestCockroachDBTestSuite(t *testing.T) {
	db := cockroach.NewCockroachDBGraph("postgresql://root@localhost:26257/linkgraph?sslmode=disable")
	c := cockroachDBTestSuite{
		cdb: db,
	}
	suite.Run(t, &c)
}

var googleLinks = []graph.Link{
	{URL: "mail.google.com"},
	{URL: "chat.google.com"},
}

var fbLinks = []graph.Link{
	{URL: "instagram.fb.com"},
	{URL: "whatsapp.fb.com"},
}

func (c *cockroachDBTestSuite) TestUpsertLink() {

	for _, link := range fbLinks {
		err := c.cdb.UpsertLink(&link)
		assert.Nil(c.T(), err)
	}
}

//func (c *cockroachDBTestSuite) TestUpsertEdge() {
//	for _, link := range googleLinks {
//		err := c.cdb.UpsertLink(&link)
//		assert.Nil(c.T(), err)
//	}
//
//	edgeOne := graph.Edge{
//		Source:      googleLinks[0].ID,
//		Destination: googleLinks[1].ID,
//	}
//	err := c.cdb.UpsertEdge(&edgeOne)
//	assert.Nil(c.T(), err)
//}

package cockroach_test

import (
	"log"
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
	db, err := cockroach.NewCockroachDbGraph("postgresql://root@localhost:26257/linkgraph?sslmode=disable")

	if err != nil{
		log.Fatalln(err)
	}

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

func (c *cockroachDBTestSuite) TestUpsertEdge() {
	for idx := range googleLinks {
		err := c.cdb.UpsertLink(&googleLinks[idx])
		assert.Nil(c.T(), err)
	}

	edgeOne := graph.Edge{
		Src: googleLinks[0].ID,
		Dst: googleLinks[1].ID,
	}

	err := c.cdb.UpsertEdge(&edgeOne)
	assert.Nil(c.T(), err)
}

package graph

import (
	"fmt"
)

var (
	ErrUnknownEdgeLinks = fmt.Errorf("unknown source/destination for edge")
	ErrNotFound = fmt.Errorf("not found")
)

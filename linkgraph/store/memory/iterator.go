package memory

import "github.com/RahilRehan/googlo/linkgraph/graph"

// compile time check to verify if the types implement interfaces
var _ graph.LinkIterator = (*linkiterator)(nil)
var _ graph.EdgeIterator = (*edgeIterator)(nil)

type linkiterator struct{
	graph *inMemoryGraph
	links []*graph.Link
	curIdx int
}

type edgeIterator struct{
	graph *inMemoryGraph
	edges []*graph.Edge
	curIdx int
}

func(i *linkiterator) Next() bool {
	if i.curIdx >= len(i.links){
		return false
	}
	i.curIdx++
	return true
}

func(i *linkiterator) Error() error{
	return nil
}

func(i *linkiterator) Close() error{
	return nil
}

func(i *linkiterator) Link() *graph.Link{
	i.graph.mu.RLock()
	defer i.graph.mu.RUnlock()

	link := new(graph.Link)
	*link = *i.links[i.curIdx-1]
	return link
}


func(i *edgeIterator) Next() bool {
	if i.curIdx >= len(i.edges){
		return false
	}
	i.curIdx++
	return true
}

func(i *edgeIterator) Error() error{
	return nil
}

func(i *edgeIterator) Close() error{
	return nil
}

func(i *edgeIterator) Edge() *graph.Edge{
	i.graph.mu.RLock()
	defer i.graph.mu.RUnlock()

	edge := new(graph.Edge)
	*edge = *i.edges[i.curIdx-1]
	return edge
}


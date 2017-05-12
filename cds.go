package fabric

// Example of a CDS: https://golang.org/src/container/list/list.go
// TODO: Convert Element to satisfy Node Interface type
//		then convert the list object to satisfy CDS Interface
//		note that elements should be wrapped to have actual integer IDs
//		that can be returned with the ID() method.

// wrap data structure elements to become generic CDS Nodes
type Node interface {
	ID() int // returns node id
}

type NodeList []Node
type EdgesMap map[Node][]Node

// add these methods to data structure objects to use as CDS
type CDS interface {
	ListNodes() NodeList
	ListEdges() EdgesMap
}

/*
	DS = Data Structure; used when a UI will have access to entire CDS
*/
type DS struct {
	Nodes []int
	Edges map[int][]int
}

func NewDS(nodes []int, edges map[int][]int) *DS {
	return &DS{
		Nodes: nodes,
		Edges: edges,
	}
}

func (s *DS) NodeCount() int {
	return len(s.Nodes)
}

func (s *DS) EdgeCount() int {
	var total int
	for i, v := range s.Edges {
		total += len(v)
	}
	return total
}

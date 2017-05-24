package fabric

// NOTE: a CDS satisfies the Section interface
type Section interface {
	ListNodes() NodeList
	ListEdges() EdgeList
}

/* Sub-graphs are non-disjoint collections of nodes and edges */
type Subgraph struct {
	Nodes NodeList
	Edges EdgeList
}

// NewSubgraph will grab all edges from nodes that connect to
// other nodes that are in our list.
func NewSubgraph(nodes NodeList, c CDS) *Subgraph {

	edges := make(EdgeList, 0)

	for _, n := range nodes {
		cdsEdges := c.ListEdges()
		for _, e := range cdsEdges {
			s := e.Source()
			d := e.Destination()
			if d.ID() == n.ID() && containsNode(nodes, s) {
				edges = append(edges, e)
			}
		}

	}

	return &Subgraph{
		Nodes: nodes,
		Edges: edges,
	}
}

func (s *Subgraph) ListNodes() NodeList {
	return s.Nodes
}

func (s *Subgraph) ListEdges() EdgeList {
	return s.Edges
}

/*
	Branches are all nodes and edges for a particuliar branch
	(usually of a tree graph)
	A branch is technically a sub-graph as well.
*/
type Branch struct {
	Nodes NodeList
	Edges EdgeList
}

func NewBranch(root Node, c CDS) *Branch {
	nodes := c.ListNodes()
	edges := c.ListEdges()

	// TODO: grab all children nodes recursively
	// dfs()

	return &Branch{
		Nodes: nodes,
		Edges: edges,
	}
}

/*
TODO:
func dfs(root Node, seen []Node, done []Node, c CDS) []Node {
	seen = append(seen, root)
	edges := c.ListEdges()
	// for each edge
	//	check if edge contains root as source
	// ...
	//  if edge contains root as source, call dfs on destination

}
*/

func (b *Branch) ListNodes() NodeList {
	return b.Nodes
}

func (b *Branch) ListEdges() EdgeList {
	return b.Edges
}

/*
	Partitions are only for linear CDSs
	(i.e. each node can only have at most 2 edges)
*/
type Partition struct {
	Nodes NodeList
	Edges EdgeList
}

func NewPartition(start, end Node, c CDS) *Partition {
	// TODO: adds all nodes between and including the start
	//		and end node; will also grab all edges for these
	//		nodes.
	nodes := c.ListNodes()
	edges := c.ListEdges()
	// TODO: recursive grab of nodes

	return &Partition{
		Nodes: nodes,
		Edges: edges,
	}
}

func (p *Partition) ListNodes() NodeList {
	return p.Nodes
}

func (p *Partition) ListEdges() EdgeList {
	return p.Edges
}

/* Subsets are used for generic node selection (but not generic edge selection) */
type Subset struct {
	Nodes NodeList
	Edges EdgeList
}

// NewSubset grabs all (and only all) edges that are connected
// to a node in the list of nodes supplied.
func NewSubset(nodes NodeList, c CDS) *Subset {
	cdsEdges := c.ListEdges()
	edges := make(EdgeList, 0)
	for _, n := range nodes {
		for _, e := range cdsEdges {
			if e.Source() == n || e.Destination() == n {
				if !containsEdge(edges, e) {
					edges = append(edges, e)
				}
			}
		}
	}

	return &Subset{
		Nodes: nodes,
		Edges: edges,
	}
}

func (s *Subset) ListNodes() NodeList {
	return s.Nodes
}

func (s *Subset) ListEdges() EdgeList {
	return s.Edges
}

/* Disjoints are a collection of arbitrary nodes and arbitrary edges */
type Disjoint struct {
	Nodes NodeList
	Edges EdgeList
}

func NewDisjoint(nodes NodeList, edges EdgeList) *Disjoint {
	return &Disjoint{
		Nodes: nodes,
		Edges: edges,
	}
}

// ComposeSections takes a list of CDS graphs (sections) and composes them into a new single disjoint
func ComposeSections(graphs []*Section) *Disjoint {
	nodes := make(NodeList, 0)
	edges := make(EdgeList, 0)

	for _, g := range graphs {
		gn := g.ListNodes()
		ge := g.ListEdges()

		// add graph nodes to Disjoint node list
		for _, n := range gn {
			if !containsNode(nodes, n) {
				nodes = append(nodes, n)
			}
		}

		// add graph edges to disjoint edge list
		for _, e := range ge {
			if !containsEdge(edges, e) {
				edges = append(edges, e)
			}
		}

	}

	return &Disjoint{
		Nodes: nodes,
		Edges: edges,
	}
}

func (d *Disjoint) ListNodes() NodeList {
	return d.Nodes
}

func (d *Disjoint) ListEdges() EdgeList {
	return d.Edges
}

package fabric

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

// Signal defines the possible signal values one dependency graph node can send to another
// RECOMMENDATION: every thread should have a proper reaction (which may be a non-reaction)
//  to each signal value for each access procedure in each dependency node.
// EXAMPLE: an example reaction to an abort signal could be "Abort Chain/Tree" where the dependents
// 	and their dependents, etc. all abort their operations if a signal value from a dependency node
// 	is an 'Abort' signal.
type Signal int

const (
	// Waiting can be used for an access procedure that has not begun but is in line too
	Waiting Signal = iota
	// Started can be used for an access procedure that is no longer waiting and has begun execution
	Started
	// Completed can be used for an access procedure that has finished execution successfully
	Completed
	// Aborted can be used for an access procedure that failed to finish execution
	Aborted
	// AbortRetry EXAMPLE: could use exponential backoff checks on retries for AbortRetry signals from dependencies ...
	AbortRetry
	// PartialAbort can be used to specify if an operation partially-completed before aborting)
	PartialAbort
)

// ProcedureSignals is used to map a signal to the access type that caused the signal
// NOTE: the string key should be equivalent to the Class() method return value for that AccessType
// EXAMPLE: a system design calls for a single thread having multiple access procedures,
// 	only some of which induce a dependent to invoke a responsive operation, then to know which
// 	procedure a signal is from you can use this map.
type ProcedureSignals map[string]Signal

// NodeType defines the possible values for types of dependency graph nodes
type NodeType int

const (
	// UINode are the spatial definitions usually assigned to a single thread
	UINode NodeType = iota
	// TemporalNode are assigned to threads which address the same UI as the temporals UI dependent
	TemporalNode
	// VirtualTemporalNode is a spawned temporary temporal node
	VirtualTemporalNode
	// VUINode is a temporary UI node
	VUINode
	// VDGNode is a node in a virtual dependency graph
	VDGNode
	// Unknown is a catch-all for an improperly constructed dependency graph node
	Unknown
)

// SignalingMap is a map of dependent node ids to a set of
// access procedures and their current signal states.
type SignalingMap map[int]chan ProcedureSignals

// SignalsMap is a map of dependency node ids to a set of
// their access procedures and their current signal states.
type SignalsMap map[int]<-chan ProcedureSignals

// DGNode (Dependency Graph Node) ...
// every DGNode has an id, a Type, a state, and a set of Access Procedures
// NOTE: This will require assigning signals to their appropriate nodes
//		when setting up a dependency graph.
type DGNode interface {
	ID() int           // must be unique from all other DGNodes in our graph
	GetType() NodeType // specifies whether node is UI, VUI, etc.
	GetPriority() int  // not necessary, but can be useful
	ListProcedures() ProcedureList
	ListDependents() []DGNode
	ListDependencies() []DGNode
	UpdateSignaling(SignalingMap, SignalsMap) // makes it possible to update the SignalingMap and SignalsMap for a DGNode
	ListSignalers() SignalingMap
	ListSignals() SignalsMap
	Signal(ProcedureSignals) // used to send the same signal to all dependents in signalers list
}

// Graph can be either UI DDAG, Temporal DAG or VDG
type Graph struct {
	DS  *CDS // reference to data structure that the dependency graph is for
	Top map[DGNode][]*DGNode
}

// NewGraph creates a new empty graph
func NewGraph() *Graph {
	return &Graph{
		Top: make(map[DGNode][]*DGNode),
	}
}

// GenID ...
func (g *Graph) GenID() int {
	rand.Seed(time.Now().UnixNano())
	id := rand.Int()
	for n := range g.Top {
		if n.ID() == id {
			id = g.GenID()
		}
	}
	return id
}

// IsLeafBoundary ...
func (g *Graph) IsLeafBoundary(n *DGNode) bool {
	if len(g.Dependents(n)) == 0 {
		return true
	}

	return false
}

// IsRootBoundary ...
func (g *Graph) IsRootBoundary(n *DGNode) bool {
	if len(g.Dependencies(n)) == 0 {
		return true
	}

	return false
}

// SignalsAndSignalers will udpate the SignalingMaps and SignalsMaps for all DGNodes in the graph
func (g *Graph) SignalsAndSignalers() {

	// for all nodes in the graph
	for n, l := range g.Top {
		// create its SignalersMap
		sm := make(SignalingMap)
		deps := g.Dependents(&n)
		for _, d := range deps {
			c := make(chan ProcedureSignals)
			sm[d.ID()] = c
		}

		// create its SignalsMap
		s := make(SignalsMap)
		for _, np := range l {
			dep := *np
			channels := dep.ListSignalers()
			ch := channels[dep.ID()]
			s[dep.ID()] = ch
		}

		n.UpdateSignaling(sm, s)
	}
}

// AddRealNode ...
// This should only be used for adding nodes to a graph
// to intialize the graph.
func (g *Graph) AddRealNode(node DGNode) error {
	if !reflect.ValueOf(node).Type().Comparable() {
		return fmt.Errorf("Node type is not comparable and cannot be used in the graph topology")
	}

	if _, ok := g.Top[node]; !ok {
		g.Top[node] = []*DGNode{}
	} else {
		return fmt.Errorf("Node already exists in Dependency Graph.")
	}
	return nil
}

// AddRealEdge will create an edge and an appropriate signaling channel between nodes
func (g *Graph) AddRealEdge(source int, dest *DGNode) {
	d := *dest

	for i, k := range g.Top {
		if i.ID() == source {
			if !containsDGNode(k, dest) {
				k = append(k, dest)

				// update SignalingMap for destination
				depSig := d.ListSignalers()
				depS := d.ListSignals()
				depSig[i.ID()] = make(chan ProcedureSignals)
				d.UpdateSignaling(depSig, depS)

				// update SignalsMap for source
				signals := i.ListSignals()
				signalers := i.ListSignalers()
				for j, v := range d.ListSignalers() {
					if j == i.ID() {
						signals[d.ID()] = v
						break
					}
				}
				i.UpdateSignaling(signalers, signals)
			}
		}
	}

}

// CycleDetect will check whether a graph has cycles or not
func (g *Graph) CycleDetect() bool {
	var seen []DGNode
	var done []DGNode

	for i := range g.Top {
		if !contains(done, i) {
			result, d := g.cycleDfs(i, seen, done)
			done = d
			if result {
				return true
			}
		}
	}
	return false
}

// Recursive Depth-First-Search; used for Cycle Detection
func (g *Graph) cycleDfs(start DGNode, seen, done []DGNode) (bool, []DGNode) {
	seen = append(seen, start)
	adj := g.Top[start]
	for _, vp := range adj {
		v := *vp
		if contains(done, v) {
			continue
		}

		if contains(seen, v) {
			return true, done
		}

		if result, done := g.cycleDfs(v, seen, done); result {
			return true, done
		}
	}
	seen = seen[:len(seen)-1]
	done = append(done, start)
	return false, done
}

// GetAdjacents will return the list of nodes a supplied node points too
func (g *Graph) GetAdjacents(node DGNode) []DGNode {
	var list []DGNode

	for n, l := range g.Top {
		if n.ID() == node.ID() {
			// Add all dependents to list
			for n2, l2 := range g.Top {
				if containsDGNode(l2, &n) {
					list = append(list, n2)
				}
			}
			// Add all dependencies to list
			for _, np := range l {
				list = append(list, *np)
			}
		}
	}

	return list
}

// TotalityUnique is a Totality-Uniqueness check for the UI nodes of a graph...
// should only be called once when creating the UI dependency graph;
// can be called with the creation of each UI if needed for
// more "real-time" verification.
func (g *Graph) TotalityUnique() bool {
	// grab all UI nodes
	var uiSlice []DGNode
	for i := range g.Top {
		if i.GetType() == UINode {
			uiSlice = append(uiSlice, i)
		}
	}

	var done []DGNode

	// for every UI Node
	for i, n := range uiSlice {
		// compare it against every other UI node
		for j, n2 := range uiSlice {
			if !contains(done, n2) {
				if j != i {
					if reflect.DeepEqual(n, n2) {
						return false
					}
				}
			}
		}
		done = append(done, n)
	}

	return true
}

// Covered returns true if all CDS nodes and edges are covered
func (g *Graph) Covered() bool {
	// grab all UI nodes
	var uiSlice []UI
	for v := range g.Top {
		if v.GetType() == UINode {
			uiSlice = append(uiSlice, v.(UI))
		}
	}

	// grab all CDS nodes and edges
	ds := *g.DS
	nodes := ds.ListNodes()
	edges := ds.ListEdges()

FIRST:
	// for every node in the CDS
	for _, v := range nodes {
		// check that at least one UI contains it
		for _, u := range uiSlice {
			s := u.GetSection()
			sp := *s
			uiCDSNodes := sp.ListNodes()
			// if UI contains node; check next CDS node
			if containsNode(uiCDSNodes, v) {
				continue FIRST
			}
		}

		// if CDS node is checked in every UI and does not show up
		return false
	}

SECOND:
	// for every edge in the CDS
	for _, v := range edges {
		// check that at least one UI contains it
		for _, u := range uiSlice {
			s := u.GetSection()
			sp := *s
			uiCDSEdges := sp.ListEdges()
			// if UI contains edge; check next CDS edge
			if containsEdge(uiCDSEdges, v) {
				continue SECOND
			}
		}

		// if CDS edge is checked in every UI and does not show up
		return false
	}

	return true
}

// AddVUI requires that the node return a true value for its IsVirtual method
func (g *Graph) AddVUI(node UI) error {
	if !node.IsVirtual() {
		return fmt.Errorf("Not a virtual node.")
	}

	var nodeSlice []DGNode
	for n := range g.Top {
		nodeSlice = append(nodeSlice, n)
	}

	if !contains(nodeSlice, node) {
		g.Top[node.(DGNode)] = []*DGNode{}
	} else {
		return fmt.Errorf("Node already exists in Dependency Graph")
	}

	return nil
}

// RemoveVUI ...
func (g *Graph) RemoveVUI(np *DGNode) error {
	n := *np
	node, ok := n.(UI)
	if !ok {
		return fmt.Errorf("Not a UI node")
	}

	if !node.IsVirtual() {
		return fmt.Errorf("Not a virtual node")
	}

	if len(node.ListDependencies()) != 0 {
		return fmt.Errorf("VUI node still has dependencies")
	}

	// Remove VUI from Signals maps in depedent nodes
	for n, l := range g.Top {
		if containsDGNode(l, np) {
			signals := n.ListSignals()
			delete(signals, node.ID())
			n.UpdateSignaling(n.ListSignalers(), signals)
		}
	}

	// remove node from graph
	delete(g.Top, node.(DGNode))

	return nil
}

// Dependents ...
func (g *Graph) Dependents(np *DGNode) []DGNode {
	var list []DGNode
	n := *np

	for i, v := range g.Top {
		if i.ID() != n.ID() {
			if containsDGNode(v, np) {
				list = append(list, i)
			}
		}
	}

	return list
}

// Dependencies ...
func (g *Graph) Dependencies(np *DGNode) []DGNode {
	var list []DGNode

	n := *np
	v, ok := g.Top[n]
	if !ok {
		return list
	}

	for _, p := range v {
		pp := *p
		list = append(list, pp)
	}
	return list
}

// Type will return the proper NodeType value for a given DGNode argument
func (g *Graph) Type(n DGNode) NodeType {
	if j, ok := n.(UI); ok {
		if j.IsVirtual() {
			return VUINode
		}
		return UINode
	}
	if k, ok := n.(Temporal); ok {
		if k.IsVirtual() {
			return VirtualTemporalNode
		}
		return TemporalNode
	}
	if _, ok := n.(Virtual); ok {
		return VDGNode
	}

	return Unknown
}

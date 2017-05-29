package ring

import (
	"github.com/JKhawaja/fabric"
)

// TODO: Implement the appropriate methods for Temporal Nodes
//		which is all DGNode methods + GetRoot() and IsVirtual()
type RingTemporal struct {
	DGNode
	UIRoot  RingUI
	Virtual bool
}

// NOTE: Our arguments ...
//		- Graph that node will be a part of
//		- UI that will be root node for temporal DAG
//		- Virtual boolean (specifying whether node is a virtual spawner or not)
//		- List of usable access procedures
// AUTOGENERATED: function signature can be autogenerated as well ...
func NewRingTemporal(g *fabric.Graph, ui RingUI, v bool, pl fabric.ProcedureList) (*RingTemporal, error) {

	// FIRST: Get ID (init)
	var R RingTemporal // Auto-generated
	// `id := g.GenID()` can be autogenerated as well
	R.Id = g.GenID()
	R.Type = g.Type(R)
	R.AccessProcedures = pl

	// SECOND: Add to graph
	err := g.AddRealNode(R) // Auto-generated
	if err != nil {
		return &R, err
	}

	// THIRD: Set up data
	R.Signalers = g.CreateSignalers(R)
	R.Signals = g.Signals(R)
	R.IsRoot = g.IsRootBoundary(R)
	R.IsLeaf = g.IsLeafBoundary(R)
	R.Dependents = g.Dependents(R)
	R.Dependencies = g.Dependencies(R)
	R.UIRoot = ui
	R.Virtual = v

	return &R, nil
}

func (r RingTemporal) ID() int {
	return r.Id
}

func (r RingTemporal) GetType() fabric.NodeType {
	return r.Type
}

func (r RingTemporal) GetPriority() int {
	// EXAMPLE: could calculate priority based on
	// priorities assigned to procedures in procedures list ...
	// Priority is used to
	p := len(r.AccessProcedures)
	return p
}

func (r RingTemporal) ListProcedures() fabric.ProcedureList {
	return r.AccessProcedures
}

func (r RingTemporal) ListDependents() []fabric.DGNode {
	return r.Dependents
}

func (r RingTemporal) ListDependencies() []fabric.DGNode {
	return r.Dependencies
}

func (r RingTemporal) ListSignals() fabric.SignalsMap {
	return r.Signals
}

func (r RingTemporal) ListSignalers() fabric.SignalingMap {
	return r.Signalers
}

func (r RingTemporal) Signal(s fabric.Signal) {
	for _, c := range r.Signalers {
		c <- s
	}
}

func (r RingTemporal) GetRoot() fabric.UI {
	return r.UIRoot
}

func (r RingTemporal) IsVirtual() bool {
	return r.Virtual
}
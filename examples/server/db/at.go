package db

import (
	"github.com/JKhawaja/fabric"
)

// AddTreeNode ...
type AddTreeNode func(*Tree, interface{}) *TreeNode

// ID ...
func (a *AddTreeNode) ID() int {
	return 0
}

// Priority ...
func (a *AddTreeNode) Priority() int {
	return 0
}

// Commit ...
func (a *AddTreeNode) Commit(np *fabric.DGNode) error {
	n := *np
	n.Signal()
	return nil
}

// Rollback ...
func (a *AddTreeNode) Rollback(np fabric.RestoreNodes, el fabric.RestoreEdges) error {
	return nil
}

// AddTreeEdge ...
type AddTreeEdge func(*Tree, *TreeNode, *TreeNode) *TreeEdge

// ID ...
func (a *AddTreeEdge) ID() int {
	return 0
}

// Priority ...
func (a *AddTreeEdge) Priority() int {
	return 0
}

// Commit ...
func (a *AddTreeEdge) Commit(np *fabric.DGNode) error {
	return nil
}

// Rollback ...
func (a *AddTreeEdge) Rollback(np fabric.RestoreNodes, el fabric.RestoreEdges) error {
	return nil
}

// DeleteTreeEntity ...
type DeleteTreeEntity func(*Tree, int)

// ID ...
func (d *DeleteTreeEntity) ID() int {
	return 0
}

// Priority ...
func (d *DeleteTreeEntity) Priority() int {
	return 0
}

// Commit ...
func (d *DeleteTreeEntity) Commit(np *fabric.DGNode) error {
	return nil
}

// Rollback ...
func (d *DeleteTreeEntity) Rollback(np fabric.RestoreNodes, el fabric.RestoreEdges) error {
	return nil
}

// ReadTreeNode ...
type ReadTreeNode func(*Tree, int) interface{}

// ID ...
func (r *ReadTreeNode) ID() int {
	return 0
}

// Priority ...
func (r *ReadTreeNode) Priority() int {
	return 0
}

// Commit ...
func (r *ReadTreeNode) Commit(np *fabric.DGNode) error {
	return nil
}

// Rollback ...
func (r *ReadTreeNode) Rollback(np fabric.RestoreNodes, el fabric.RestoreEdges) error {
	return nil
}

// UpdateTreeNode ...
type UpdateTreeNode func(*Tree, int, interface{})

// ID ...
func (u *UpdateTreeNode) ID() int {
	return 0
}

// Priority ...
func (u *UpdateTreeNode) Priority() int {
	return 0
}

// Commit ...
func (u *UpdateTreeNode) Commit(np *fabric.DGNode) error {
	return nil
}

// Rollback ...
func (u *UpdateTreeNode) Rollback(np fabric.RestoreNodes, el fabric.RestoreEdges) error {
	return nil
}
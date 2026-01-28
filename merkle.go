package merkle

import (
	"bytes"

	"github.com/ethereum/go-ethereum/crypto"
)

type Node struct {
	Hash       []byte
	MerkleHash []byte
	Priority   uint64
	Left       *Node
	Right      *Node
}

type Treap struct {
	Root *Node
}

type ITreap interface {
	Insert(key []byte, priority uint64)
	Search(key []byte) *Node
	Delete(key []byte)
	GetRootHash() []byte
}

var _ ITreap = &Treap{}

func NewTreap() *Treap {
	return &Treap{}
}

func (t *Treap) GetRootHash() []byte {
	if t.Root == nil {
		return nil
	}
	return t.Root.MerkleHash
}

func (t *Treap) Insert(key []byte, priority uint64) {
	node := &Node{
		Hash:     key,
		Priority: priority,
	}

	if t.Root == nil {
		t.Root = node
		return
	}

	leftSubtree, rightSubtree := split(t.Root, key)
	t.Root = merge(merge(leftSubtree, node), rightSubtree)
}

func (t *Treap) Search(key []byte) *Node {
	cur := t.Root

	for cur != nil {
		cmp := bytes.Compare(key, cur.Hash)

		if cmp == 0 {
			return cur
		} else if cmp < 0 {
			cur = cur.Left
		} else {
			cur = cur.Right
		}
	}

	return nil
}

func (t *Treap) Delete(key []byte) {
	if t.Root == nil {
		return
	}

	leftSubtree, rightSubtree := split(t.Root, key)

	if bytes.Equal(rightSubtree.Hash, key) {
		t.Root = merge(leftSubtree, rightSubtree.Right)
		return
	}

	var stack []*Node
	cur := rightSubtree
	for cur != nil && cur.Left != nil {
		if bytes.Equal(cur.Left.Hash, key) {
			cur.Left = merge(cur.Left.Left, cur.Left.Right)
			updateNode(cur)
			for i := len(stack) - 1; i >= 0; i-- {
				updateNode(stack[i])
			}
			break
		}
		stack = append(stack, cur)
		cur = cur.Left
	}

	t.Root = merge(leftSubtree, rightSubtree)
}

func split(root *Node, key []byte) (*Node, *Node) {
	if root == nil {
		return nil, nil
	}

	if bytes.Compare(root.Hash, key) < 0 {
		leftSubtree, rightSubtree := split(root.Right, key)
		root.Right = leftSubtree
		updateNode(root)
		return root, rightSubtree
	} else {
		leftSubtree, rightSubtree := split(root.Left, key)
		root.Left = rightSubtree
		updateNode(root)
		return leftSubtree, root
	}
}

func merge(left, right *Node) *Node {
	if left == nil {
		return right
	}
	if right == nil {
		return left
	}

	if left.Priority > right.Priority {
		left.Right = merge(left.Right, right)
		updateNode(left)
		return left
	}

	right.Left = merge(left, right.Left)
	updateNode(right)
	return right
}

func updateNode(node *Node) {
	if node == nil {
		return
	}
	childrenHash := hashNodes(node.Left, node.Right)
	if childrenHash == nil {
		node.MerkleHash = node.Hash
	} else {
		node.MerkleHash = hash(childrenHash, node.Hash)
	}
}

func hashNodes(left, right *Node) []byte {
	var lHash, rHash []byte

	if left != nil {
		lHash = left.MerkleHash
	}
	if right != nil {
		rHash = right.MerkleHash
	}

	return hash(lHash, rHash)
}

func hash(a, b []byte) []byte {
	if len(a) == 0 {
		return b
	}
	if len(b) == 0 {
		return a
	}

	if bytes.Compare(a, b) < 0 {
		return crypto.Keccak256(a, b)
	}
	return crypto.Keccak256(b, a)
}

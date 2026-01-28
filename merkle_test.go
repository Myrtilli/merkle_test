package merkle

import (
	"bytes"
	"testing"
)

func TestTreap_InsertSearchDelete(t *testing.T) {
	tr := NewTreap()

	keys := [][]byte{
		[]byte("a"),
		[]byte("b"),
		[]byte("c"),
		[]byte("d"),
	}
	priorities := []uint64{10, 20, 5, 15}

	for i, key := range keys {
		tr.Insert(key, priorities[i])
		if tr.Search(key) == nil {
			t.Fatalf("Insert failed: key %s not found", key)
		}
	}

	rootHashAfterInsert := tr.GetRootHash()
	if rootHashAfterInsert == nil {
		t.Fatal("Root hash is nil after inserts")
	}

	for _, key := range keys {
		node := tr.Search(key)
		if node == nil {
			t.Fatalf("Search failed: key %s not found", key)
		}
		if !bytes.Equal(node.Hash, key) {
			t.Fatalf("Search returned wrong node: got %s, want %s", node.Hash, key)
		}
	}

	if tr.Search([]byte("x")) != nil {
		t.Fatal("Search should return nil for non-existent key")
	}

	tr.Delete([]byte("b"))
	if tr.Search([]byte("b")) != nil {
		t.Fatal("Delete failed: key b still found")
	}

	tr.Delete([]byte("a"))
	if tr.Search([]byte("a")) != nil {
		t.Fatal("Delete failed: key a still found")
	}

	tr.Delete([]byte("d"))
	if tr.Search([]byte("d")) != nil {
		t.Fatal("Delete failed: key d still found")
	}

	node := tr.Search([]byte("c"))
	if node == nil || !bytes.Equal(node.Hash, []byte("c")) {
		t.Fatal("After deletes, remaining key c not found")
	}

	rootHashAfterDelete := tr.GetRootHash()
	if bytes.Equal(rootHashAfterInsert, rootHashAfterDelete) {
		t.Fatal("Root hash should change after deletions")
	}
}

package whisker

import (
    "errors"
    "fmt"

//    "github.com/Aanthord/remy-go/pkg/dna"
    "github.com/Aanthord/remy-go/pkg/memory"
)

// WhiskerTree represents a tree-like structure that holds multiple whiskers
type WhiskerTree struct {
    Root *WhiskerNode
}

// WhiskerNode represents a node in the whisker tree
type WhiskerNode struct {
    Whisker  *Whisker
    Children []*WhiskerNode
}

// NewWhiskerTree is a constructor that creates a new instance of the WhiskerTree struct
func NewWhiskerTree() *WhiskerTree {
    root := &WhiskerNode{
        Whisker: NewWhisker(0, 0, 1.0, 0.0, memory.NewMemoryRange(memory.MinMemory(), memory.MaxMemory())),
    }
    return &WhiskerTree{Root: root}
}

// Insert inserts a new whisker into the whisker tree
func (wt *WhiskerTree) Insert(whisker *Whisker) error {
    if err := wt.insert(wt.Root, whisker); err != nil {
        return fmt.Errorf("failed to insert whisker: %v", err)
    }
    return nil
}

// insert is a recursive function that inserts a new whisker into the whisker tree
func (wt *WhiskerTree) insert(node *WhiskerNode, whisker *Whisker) error {
    if node.Whisker.Domain.Intersects(whisker.Domain) {
        if node.Whisker.Generation >= whisker.Generation {
            return fmt.Errorf("whisker with generation %d already exists in the domain", whisker.Generation)
        }
        node.Whisker = whisker
        return nil
    }

    for _, child := range node.Children {
        if child.Whisker.Domain.Intersects(whisker.Domain) {
            return wt.insert(child, whisker)
        }
    }

    newNode := &WhiskerNode{Whisker: whisker}
    node.Children = append(node.Children, newNode)
    return nil
}

// FindWhisker finds the whisker that corresponds to the given memory state
func (wt *WhiskerTree) FindWhisker(m *memory.Memory) (*Whisker, error) {
    node, err := wt.findNode(wt.Root, m)
    if err != nil {
        return nil, err
    }
    return node.Whisker, nil
}

// findNode is a recursive function that finds the node that contains the given memory state
func (wt *WhiskerTree) findNode(node *WhiskerNode, m *memory.Memory) (*WhiskerNode, error) {
    if node.Whisker.Domain.Contains(m) {
        return node, nil
    }

    for _, child := range node.Children {
        if child.Whisker.Domain.Contains(m) {
            return wt.findNode(child, m)
        }
    }

    return nil, errors.New("memory not found in the tree")
}

// String returns a string representation of the whisker tree
func (wt *WhiskerTree) String() string {
    return wt.toString(wt.Root, 0)
}

// toString is a recursive function that generates a string representation of the whisker tree
func (wt *WhiskerTree) toString(node *WhiskerNode, indent int) string {
    str := fmt.Sprintf("%s%s\n", getIndent(indent), node.Whisker)
    for _, child := range node.Children {
        str += wt.toString(child, indent+2)
    }
    return str
}

// getIndent is a helper function that returns a string of spaces based on the given indentation level
func getIndent(indent int) string {
    return fmt.Sprintf("%*s", indent, "")
}

package lunar

// Node is a tree node
type Node struct {
	Name     string
	Value    string // value only makes sense in leaf
	Children []*Node
}

// AddChildren adds children node if it's not existing
func (node *Node) AddChildren(children ...*Node) {
	for _, child := range children {
		if child != nil && node.GetChild(child.Name) == nil {
			node.Children = append(node.Children, child)
		}
	}
}

// GetChild gets a child by its name
func (node *Node) GetChild(name string) *Node {
	for _, child := range node.Children {
		if child.Name == name {
			return child
		}
	}

	return nil
}

// IsLeaf returns true if the node is a leaf
func (node *Node) IsLeaf() bool {
	return len(node.Children) == 0
}

// ToMap converts to map
func (node *Node) ToMap() interface{} {
	if node.IsLeaf() {
		return node.Value
	}

	m := make(map[string]interface{})
	for _, child := range node.Children {
		m[child.Name] = child.ToMap()
	}

	return m
}

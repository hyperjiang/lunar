package lunar

import (
	"encoding/json"
	"strings"
)

// Items is items under namespace
type Items map[string]string

// Get gets value of given key
func (items Items) Get(key string) string {
	if v, ok := items[key]; ok {
		return v
	}

	return ""
}

// String converts Items to json string
func (items Items) String() string {
	bytes, _ := json.Marshal(items.parse())
	return string(bytes)
}

func (items Items) parse() interface{} {
	var root Node
	for k, v := range items {
		ks := strings.Split(k, ".")

		parent := root.GetChild(ks[0])
		if parent == nil {
			parent = &Node{
				Name:  ks[0],
				Value: v,
			}
			root.AddChildren(parent)
		}

		node := parent
		for i := 1; i < len(ks); i++ {
			child := node.GetChild(ks[i])
			if child == nil {
				child = &Node{
					Name:  ks[i],
					Value: v,
				}
				node.AddChildren(child)
			}
			node = child
		}
	}

	return root.ToMap()
}

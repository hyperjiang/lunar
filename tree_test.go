package lunar

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTree(t *testing.T) {
	n0 := &Node{
		Name: "root",
	}
	n1 := &Node{
		Name: "a",
	}
	n2 := &Node{
		Name:  "b",
		Value: "b",
	}
	n3 := &Node{
		Name:  "c",
		Value: "c",
	}
	n4 := &Node{
		Name:  "d",
		Value: "d",
	}
	n5 := &Node{
		Name:  "e",
		Value: "e",
	}

	n1.AddChildren(n2, n3)
	n4.AddChildren(n5)
	n0.AddChildren(n1, n4)

	m := map[string]interface{}{
		"a": map[string]interface{}{
			"b": "b",
			"c": "c",
		},
		"d": map[string]interface{}{
			"e": "e",
		},
	}

	should := require.New(t)
	should.Equal(m, n0.ToMap())
}

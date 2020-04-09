package lunar

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeURL(t *testing.T) {
	should := require.New(t)

	tests := []struct {
		url  string
		want string
	}{
		{
			"localhost",
			"http://localhost",
		},
		{
			"http://localhost",
			"http://localhost",
		},
		{
			"https://localhost",
			"https://localhost",
		},
	}

	for _, test := range tests {
		should.Equal(test.want, normalizeURL(test.url))
	}
}

func TestGetFormat(t *testing.T) {
	should := require.New(t)

	tests := []struct {
		str  string
		want string
	}{
		{
			"application",
			"properties",
		},
		{
			"application.properties",
			"properties",
		},
		{
			"abc.json",
			"json",
		},
		{
			"abc.def.yaml",
			"yaml",
		},
		{
			"abc.common",
			"properties",
		},
	}

	for _, test := range tests {
		should.Equal(test.want, GetFormat(test.str))
	}
}

func TestExpand(t *testing.T) {
	should := require.New(t)

	should.Nil(expand([]string{}, ""))

	m := make(map[string]interface{})
	m["a"] = "1"
	should.Equal(m, expand([]string{"a"}, "1"))
	should.Equal(m, Expand("a", "1"))

	m2 := make(map[string]interface{})
	m2["b"] = m
	should.Equal(m2, expand([]string{"b", "a"}, "1"))
	should.Equal(m2, Expand("b.a", "1"))

	m3 := make(map[string]interface{})
	m3["c"] = m2
	should.Equal(m3, expand([]string{"c", "b", "a"}, "1"))
	should.Equal(m3, Expand("c.b.a", "1"))

	m4 := make(map[string]interface{})
	m4["d"] = m3
	should.Equal(m4, expand([]string{"d", "c", "b", "a"}, "1"))
	should.Equal(m4, Expand("d.c.b.a", "1"))
}

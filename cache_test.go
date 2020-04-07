package lunar

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CacheTestSuite struct {
	suite.Suite
}

// TestCacheTestSuite runs the Cache test suite
func TestCacheTestSuite(t *testing.T) {
	suite.Run(t, new(CacheTestSuite))
}

func (ts *CacheTestSuite) TestMemoryCache() {
	should := require.New(ts.T())

	items := make(Items)
	items["a"] = "apple"
	items["b"] = "banana"

	items2 := make(Items)
	items2["content"] = "this is plaintext"

	cache := new(MemoryCache)
	cache.SetItems("ns", items)
	cache.SetItems("ns.txt", items2)

	should.Equal("apple", cache.GetItems("ns").Get("a"))
	should.Equal("this is plaintext", cache.GetItems("ns.txt").Get("content"))
	should.ElementsMatch([]string{"ns", "ns.txt"}, cache.GetKeys())

	cache.Delete("ns")
	should.Len(cache.GetItems("ns"), 0)
	should.Len(cache.GetKeys(), 1)

	cache.Drain()
	should.Len(cache.GetKeys(), 0)
}

func (ts *CacheTestSuite) TestFileCache() {
	should := require.New(ts.T())

	items := make(Items)
	items["a"] = "apple"
	items["b"] = "banana"

	items2 := make(Items)
	items2["content"] = "this is plaintext"

	cache := NewFileCache("myApp", "/tmp")
	cache.SetItems("ns", items)
	cache.SetItems("ns.txt", items2)

	should.Equal("apple", cache.GetItems("ns").Get("a"))
	should.Equal("this is plaintext", cache.GetItems("ns.txt").Get("content"))
	should.ElementsMatch([]string{"ns", "ns.txt"}, cache.GetKeys())

	cache.Delete("ns")
	should.Len(cache.GetItems("ns"), 0)
	should.Len(cache.GetKeys(), 1)

	cache.Drain()
	should.Len(cache.GetKeys(), 0)
}

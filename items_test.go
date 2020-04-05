package lunar

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ItemsTestSuite struct {
	suite.Suite
}

// TestItemsTestSuite runs the Items test suite
func TestItemsTestSuite(t *testing.T) {
	suite.Run(t, new(ItemsTestSuite))
}

func (ts *ItemsTestSuite) TestGet() {
	should := require.New(ts.T())

	items := make(Items)
	should.Empty(items.Get("foo"))

	items["foo"] = "bar"
	should.Equal("bar", items.Get("foo"))
}

func (ts *ItemsTestSuite) TestString() {
	should := require.New(ts.T())

	items := make(Items)
	items["a"] = "apple"
	items["b"] = "banana"

	should.Equal("{\"a\":\"apple\",\"b\":\"banana\"}", items.String())
}

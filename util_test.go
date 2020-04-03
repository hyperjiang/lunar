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

func TestSplitCommaSeparatedURL(t *testing.T) {
	should := require.New(t)

	urls := splitCommaSeparatedURL("192.168.1.1, http://localhost/")

	want := []string{"http://192.168.1.1", "http://localhost"}

	should.Equal(want, urls)
}

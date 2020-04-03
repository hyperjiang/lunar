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
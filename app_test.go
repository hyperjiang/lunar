package lunar

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	gock "gopkg.in/h2non/gock.v1"
)

type LunarTestSuite struct {
	suite.Suite
	app *App
}

// TestLunarTestSuite runs the Lunar test suite
func TestLunarTestSuite(t *testing.T) {
	suite.Run(t, new(LunarTestSuite))
}

// SetupSuite run once at the very start of the testing suite, before any tests are run.
func (ts *LunarTestSuite) SetupSuite() {
	ts.app = New("SampleApp", WithLogger(Printf))
}

// TearDownSuite run once at the very end of the testing suite, after all tests have been run.
func (ts *LunarTestSuite) TearDownSuite() {
	gock.Off()
}

func (ts *LunarTestSuite) mockGetNamespace(namespace string) {
	resBody, _ := ioutil.ReadFile("./mocks/GetNamespace_" + namespace + ".json")

	url := fmt.Sprintf("/configs/%s/%s/%s",
		url.QueryEscape(ts.app.ID),
		url.QueryEscape(ts.app.Cluster),
		url.QueryEscape(namespace),
	)

	gock.New(ts.app.Server).
		Get(url).
		Reply(http.StatusOK).
		BodyString(string(resBody))
}

func (ts *LunarTestSuite) TestGetValue() {
	var should = require.New(ts.T())

	ts.mockGetNamespace(defaultNamespace)
	v, err := ts.app.GetValue("portal.elastic.document.type")

	should.NoError(err)
	should.Equal("biz", v)
}

func (ts *LunarTestSuite) TestGetItems() {
	var should = require.New(ts.T())

	ts.mockGetNamespace(defaultNamespace)
	items, err := ts.app.GetItems()

	should.NoError(err)
	should.Contains(items, "portal.elastic.document.type")
}

func (ts *LunarTestSuite) TestGetContent() {
	var should = require.New(ts.T())

	ns := "a.txt"

	ts.mockGetNamespace(ns)
	content, err := ts.app.GetContent(ns)

	should.NoError(err)
	should.Equal("greetings", content)

	ns = defaultNamespace

	ts.mockGetNamespace(ns)
	content, err = ts.app.GetContent(ns)

	should.NoError(err)
	should.Equal("{\"portal.elastic.cluster.name\":\"hermes-es-fws\",\"portal.elastic.document.type\":\"biz\"}", content)
}

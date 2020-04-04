package lunar

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	gock "gopkg.in/h2non/gock.v1"
)

type ApolloClientTestSuite struct {
	suite.Suite
	client *ApolloClient
}

// TestApolloClientTestSuite runs the ApolloClient test suite
func TestApolloClientTestSuite(t *testing.T) {
	suite.Run(t, new(ApolloClientTestSuite))
}

// SetupSuite run once at the very start of the testing suite, before any tests are run.
func (ts *ApolloClientTestSuite) SetupSuite() {
	ts.client = NewApolloClient("SampleApp", WithServer("localhost:8080"), WithLogger(Printf))
}

// TearDownSuite run once at the very end of the testing suite, after all tests have been run.
func (ts *ApolloClientTestSuite) TearDownSuite() {
	gock.Off()
}

func (ts *ApolloClientTestSuite) TestGetCachedItems() {
	var should = require.New(ts.T())

	resBody, err := ioutil.ReadFile("./mocks/GetCachedItems.json")
	should.NoError(err)

	url := fmt.Sprintf("/configfiles/json/%s/%s/%s",
		url.QueryEscape(ts.client.AppID),
		url.QueryEscape(ts.client.Cluster),
		url.QueryEscape(defaultNamespace),
	)

	gock.New(ts.client.Server).
		Get(url).
		Reply(200).
		BodyString(string(resBody))

	res, err := ts.client.GetCachedItems("application")

	should.NoError(err)
	should.Contains(res, "portal.elastic.document.type")
}

func (ts *ApolloClientTestSuite) TestGetNamespace() {
	var should = require.New(ts.T())

	resBody, err := ioutil.ReadFile("./mocks/GetNamespace_application.json")
	should.NoError(err)

	url := fmt.Sprintf("/configs/%s/%s/%s",
		url.QueryEscape(ts.client.AppID),
		url.QueryEscape(ts.client.Cluster),
		url.QueryEscape(defaultNamespace),
	)

	gock.New(ts.client.Server).
		Get(url).
		Reply(200).
		BodyString(string(resBody))

	res, err := ts.client.GetNamespace("", "")

	should.NoError(err)
	should.Equal("20170430092936-dee2d58e74515ff3", res.ReleaseKey)
}

func (ts *ApolloClientTestSuite) TestGetNotifications() {
	var should = require.New(ts.T())

	resBody, err := ioutil.ReadFile("./mocks/GetNotifications.json")
	should.NoError(err)

	gock.New(ts.client.Server).
		Get("/notifications/v2").
		MatchParam("appId", ts.client.AppID).
		MatchParam("cluster", ts.client.Cluster).
		Reply(200).
		BodyString(string(resBody))

	res, err := ts.client.GetNotifications(nil)

	should.NoError(err)
	should.Equal(defaultNamespace, res[0].Namespace)
	should.Equal(101, res[0].NotificationID)
}

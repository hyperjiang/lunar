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
func (suite *ApolloClientTestSuite) SetupSuite() {
	suite.client = NewApolloClient("SampleApp", WithServer("localhost:8080"), WithLogger(Printf))
	gock.InterceptClient(suite.client.Client)
}

// TearDownSuite run once at the very end of the testing suite, after all tests have been run.
func (suite *ApolloClientTestSuite) TearDownSuite() {
	gock.RestoreClient(suite.client.Client)
	gock.Off()
}

func (suite *ApolloClientTestSuite) TestGetCachedItems() {
	var should = require.New(suite.T())

	resBody, err := ioutil.ReadFile("./mocks/GetCachedItems.json")
	should.NoError(err)

	url := fmt.Sprintf("/configfiles/json/%s/%s/%s",
		url.QueryEscape(suite.client.AppID),
		url.QueryEscape(suite.client.Cluster),
		url.QueryEscape(defaultNamespace),
	)

	gock.New(suite.client.Server).
		Get(url).
		Reply(200).
		BodyString(string(resBody))

	res, err := suite.client.GetCachedItems("application")

	should.NoError(err)
	should.Contains(res, "portal.elastic.document.type")
}

func (suite *ApolloClientTestSuite) TestGetNamespace() {
	var should = require.New(suite.T())

	resBody, err := ioutil.ReadFile("./mocks/GetNamespace.json")
	should.NoError(err)

	url := fmt.Sprintf("/configs/%s/%s/%s",
		url.QueryEscape(suite.client.AppID),
		url.QueryEscape(suite.client.Cluster),
		url.QueryEscape(defaultNamespace),
	)

	gock.New(suite.client.Server).
		Get(url).
		Reply(200).
		BodyString(string(resBody))

	res, err := suite.client.GetNamespace("", "")

	should.NoError(err)
	should.Equal("20170430092936-dee2d58e74515ff3", res.ReleaseKey)
}

func (suite *ApolloClientTestSuite) TestGetNotifications() {
	var should = require.New(suite.T())

	resBody, err := ioutil.ReadFile("./mocks/GetNotifications.json")
	should.NoError(err)

	gock.New(suite.client.Server).
		Get("/notifications/v2").
		MatchParam("appId", suite.client.AppID).
		MatchParam("cluster", suite.client.Cluster).
		Reply(200).
		BodyString(string(resBody))

	res, err := suite.client.GetNotifications(nil)

	should.NoError(err)
	should.Equal(defaultNamespace, res[0].Namespace)
	should.Equal(101, res[0].NotificationID)
}

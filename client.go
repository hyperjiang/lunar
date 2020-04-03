package lunar

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultServer         = "localhost:8080"
	defaultCluster        = "default"
	defaultNamespace      = "application"
	defaultNotificationID = -1
	defaultClientTimeout  = time.Second * 90
)

// ApolloClient is the implementation of apollo client.
//
// https://github.com/ctripcorp/apollo/wiki/%E5%85%B6%E5%AE%83%E8%AF%AD%E8%A8%80%E5%AE%A2%E6%88%B7%E7%AB%AF%E6%8E%A5%E5%85%A5%E6%8C%87%E5%8D%97
type ApolloClient struct {
	Client   *http.Client
	Server   string
	AppID    string
	Cluster  string
	ClientIP string
	logger   Logger
}

// ApolloClientOption is apollo client option
type ApolloClientOption func(*ApolloClient)

// Configurations is apollo configurations
type Configurations map[string]string

// WithServer sets apollo server address
func WithServer(server string) ApolloClientOption {
	return func(a *ApolloClient) {
		a.Server = server
	}
}

// WithAppID sets apollo app id
func WithAppID(appID string) ApolloClientOption {
	return func(a *ApolloClient) {
		a.AppID = appID
	}
}

// WithCluster sets apollo cluster
func WithCluster(cluster string) ApolloClientOption {
	return func(a *ApolloClient) {
		a.Cluster = cluster
	}
}

// WithLogger sets logger
func WithLogger(logger Logger) ApolloClientOption {
	return func(a *ApolloClient) {
		a.logger = logger
	}
}

// NewApolloClient creates a apollo client
func NewApolloClient(opts ...ApolloClientOption) (*ApolloClient, error) {
	c := &ApolloClient{}
	for _, opt := range opts {
		opt(c)
	}

	if c.AppID == "" {
		return nil, errors.New("app id can not be empty")
	}

	c.Client = &http.Client{
		Timeout: defaultClientTimeout,
	}

	c.ClientIP = getLocalIP()

	if c.Server == "" {
		c.Server = defaultServer
	}

	c.Server = normalizeURL(c.Server)

	if c.Cluster == "" {
		c.Cluster = defaultCluster
	}

	if c.logger == nil {
		c.logger = DefaultLogger
	}

	return c, nil
}

func (c *ApolloClient) get(url string, result interface{}) error {
	c.logger.Printf("%s", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		err = json.Unmarshal(body, result)
	}

	c.logger.Printf("http status: %d", resp.StatusCode)

	return err
}

// GetCachedConfigs gets cached configs from apollo
func (c *ApolloClient) GetCachedConfigs(namespace string) (Configurations, error) {
	url := fmt.Sprintf("%s/configfiles/json/%s/%s/%s?ip=%s",
		c.Server,
		url.QueryEscape(c.AppID),
		url.QueryEscape(c.Cluster),
		url.QueryEscape(namespace),
		c.ClientIP,
	)

	var res Configurations
	err := c.get(url, &res)

	return res, err
}

// GetConfigsRequest is apollo request
type GetConfigsRequest struct {
	Namespace  string
	ReleaseKey string
}

// GetConfigsResponse is apollo response
type GetConfigsResponse struct {
	AppID          string         `json:"appId"`
	Cluster        string         `json:"cluster"`
	NamespaceName  string         `json:"namespaceName"`
	Configurations Configurations `json:"configurations"`
	ReleaseKey     string         `json:"releaseKey"`
}

// GetConfigs gets realtime configs from apollo
func (c *ApolloClient) GetConfigs(r GetConfigsRequest) (*GetConfigsResponse, error) {
	if r.Namespace == "" {
		r.Namespace = defaultNamespace
	}

	url := fmt.Sprintf("%s/configs/%s/%s/%s?releaseKey=%s&ip=%s",
		c.Server,
		url.QueryEscape(c.AppID),
		url.QueryEscape(c.Cluster),
		url.QueryEscape(r.Namespace),
		url.QueryEscape(r.ReleaseKey),
		c.ClientIP,
	)

	var res GetConfigsResponse
	err := c.get(url, &res)

	return &res, err
}

// Notifications is a set of notifications
type Notifications []Notification

// String converts Notifications to json string
func (ns Notifications) String() string {
	bytes, _ := json.Marshal(ns)
	return string(bytes)
}

// Notification is the definition of notification
type Notification struct {
	Namespace      string `json:"namespaceName"`
	NotificationID int    `json:"notificationId"`
}

// GetNotifications gets notifications from apollo
func (c *ApolloClient) GetNotifications(ns Notifications) (Notifications, error) {
	if ns == nil || len(ns) == 0 {
		ns = append(ns, Notification{Namespace: defaultNamespace, NotificationID: defaultNotificationID})
	}

	url := fmt.Sprintf("%s/notifications/v2?appId=%s&cluster=%s&notifications=%s",
		c.Server,
		url.QueryEscape(c.AppID),
		url.QueryEscape(c.Cluster),
		url.QueryEscape(ns.String()),
	)

	var res Notifications
	err := c.get(url, &res)

	return res, err
}

package lunar

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// ApolloAPI is the interface of apollo api
type ApolloAPI interface {
	GetCachedItems(namespace string) (Items, error)
	GetNamespace(namespace string, releaseKey string) (*Namespace, error)
	GetNotifications(ns Notifications) (Notifications, error)
}

// ApolloClient is the implementation of apollo client.
//
// https://github.com/ctripcorp/apollo/wiki/%E5%85%B6%E5%AE%83%E8%AF%AD%E8%A8%80%E5%AE%A2%E6%88%B7%E7%AB%AF%E6%8E%A5%E5%85%A5%E6%8C%87%E5%8D%97
type ApolloClient struct {
	Options  // inherited options
	AppID    string
	Client   *http.Client
	ClientIP string
}

// make sure ApolloClient implements ApolloAPI
var _ ApolloAPI = new(ApolloClient)

// NewApolloClient creates a apollo client
func NewApolloClient(appID string, opts ...Option) *ApolloClient {
	c := &ApolloClient{
		AppID:    appID,
		Options:  NewOptions(opts...),
		ClientIP: GetLocalIP(),
	}

	c.Client = &http.Client{
		Timeout: c.ClientTimeout,
	}

	return c
}

func (c *ApolloClient) get(pathWithQuery string, result interface{}) error {
	url := c.Server + pathWithQuery
	c.Logger.Printf("%s", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	if len(c.AccessKeySecret) > 0 {
		headers := buildHeaders(pathWithQuery, c.AppID, c.AccessKeySecret)
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		err = json.Unmarshal(body, result)
	}

	c.Logger.Printf("[%d] %s", resp.StatusCode, body)

	return err
}

// GetCachedItems gets cached configs from apollo
func (c *ApolloClient) GetCachedItems(namespace string) (Items, error) {
	url := fmt.Sprintf("/configfiles/json/%s/%s/%s?ip=%s",
		url.QueryEscape(c.AppID),
		url.QueryEscape(c.Cluster),
		url.QueryEscape(namespace),
		c.ClientIP,
	)

	var res Items
	err := c.get(url, &res)

	return res, err
}

// Namespace is apollo namespace data
type Namespace struct {
	AppID      string `json:"appId"`
	Cluster    string `json:"cluster"`
	Name       string `json:"namespaceName"`
	Items      Items  `json:"configurations"`
	ReleaseKey string `json:"releaseKey"`
}

// GetNamespace gets realtime namespace data from apollo
func (c *ApolloClient) GetNamespace(namespace string, releaseKey string) (*Namespace, error) {
	if namespace == "" {
		namespace = defaultNamespace
	}

	url := fmt.Sprintf("/configs/%s/%s/%s?releaseKey=%s&ip=%s",
		url.QueryEscape(c.AppID),
		url.QueryEscape(c.Cluster),
		url.QueryEscape(namespace),
		url.QueryEscape(releaseKey),
		c.ClientIP,
	)

	var res Namespace
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
	if len(ns) == 0 {
		ns = append(ns, Notification{Namespace: defaultNamespace, NotificationID: defaultNotificationID})
	}

	url := fmt.Sprintf("/notifications/v2?appId=%s&cluster=%s&notifications=%s",
		url.QueryEscape(c.AppID),
		url.QueryEscape(c.Cluster),
		url.QueryEscape(ns.String()),
	)

	var res Notifications
	err := c.get(url, &res)

	return res, err
}

package lunar

import (
	"time"
)

const (
	defaultServer         = "localhost:8080"
	defaultCluster        = "default"
	defaultNamespace      = "application"
	defaultNotificationID = -1
	defaultClientTimeout  = time.Second * 90
)

// Options is common options
type Options struct {
	Server        string
	AppID         string
	Cluster       string
	ClientTimeout time.Duration
	Logger        Logger
}

// NewOptions creates options with defaults
func NewOptions(opts ...Option) Options {
	var options = Options{
		Server:        defaultServer,
		Cluster:       defaultCluster,
		ClientTimeout: defaultClientTimeout,
		Logger:        defaultLogger,
	}
	for _, opt := range opts {
		opt(&options)
	}

	return options
}

// Option is for setting options
type Option func(*Options)

// WithServer sets apollo server address
func WithServer(server string) Option {
	return func(o *Options) {
		o.Server = normalizeURL(server)
	}
}

// WithAppID sets apollo app id
func WithAppID(appID string) Option {
	return func(o *Options) {
		o.AppID = appID
	}
}

// WithCluster sets apollo cluster
func WithCluster(cluster string) Option {
	return func(o *Options) {
		o.Cluster = cluster
	}
}

// WithLogger sets logger
func WithLogger(logger Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

// WithClientTimeout sets client timeout
func WithClientTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.ClientTimeout = timeout
	}
}

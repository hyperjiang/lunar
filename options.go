package lunar

import (
	"strings"
	"time"
)

const (
	defaultServer           = "localhost:8080"
	defaultCluster          = "default"
	defaultNamespace        = "application"
	defaultFormat           = "properties"
	defaultNotificationID   = -1
	defaultClientTimeout    = time.Second * 90
	defaultLongPollInterval = time.Second
)

// Options is common options
type Options struct {
	Server           string
	Cluster          string
	AccessKeySecret  string
	Logger           Logger
	ClientTimeout    time.Duration
	LongPollInterval time.Duration
}

// NewOptions creates options with defaults
func NewOptions(opts ...Option) Options {
	var options = Options{
		Server:           normalizeURL(defaultServer),
		Cluster:          defaultCluster,
		ClientTimeout:    defaultClientTimeout,
		LongPollInterval: defaultLongPollInterval,
		Logger:           defaultLogger,
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

// WithLongPollInterval sets long poll interval
func WithLongPollInterval(interval time.Duration) Option {
	return func(o *Options) {
		o.LongPollInterval = interval
	}
}

// WithAccessKeySecret sets access key secret
func WithAccessKeySecret(secret string) Option {
	return func(o *Options) {
		o.AccessKeySecret = strings.TrimSpace(secret)
	}
}

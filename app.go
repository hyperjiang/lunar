package lunar

import (
	"sync"
)

// Lunar is the interface of lunar
type Lunar interface {
	GetValue(key string) (string, error)
	GetValueInNamespace(key string, namespace string) (string, error)
	GetItems() (Items, error)
	GetItemsInNamespace(namespace string) (Items, error)
	GetContent(namespace string) (string, error)
}

// App represents a single application, an application has a unique app id and manage multiple namespaces.
type App struct {
	Options           // inherited options
	ID      string    // app id
	Client  ApolloAPI // the apollo client
	meta    sync.Map  // key: namespace name, value: NamespaceMeta
}

// make sure App implements Lunar
var _ Lunar = new(App)

// NamespaceMeta is namespace metadata
type NamespaceMeta struct {
	ReleaseKey     string
	NotificationID int
}

// New creates an application, user must specify the correct app id
func New(appID string, opts ...Option) *App {
	app := &App{
		ID:      appID,
		Options: NewOptions(opts...),
	}

	app.Client = NewApolloClient(appID, opts...)

	return app
}

// UseClient sets the underlying apollo client
func (app *App) UseClient(client ApolloAPI) *App {
	app.Client = client

	return app
}

// GetValue gets value of key in default namespace
func (app *App) GetValue(key string) (string, error) {
	return app.GetValueInNamespace(key, defaultNamespace)
}

// GetValueInNamespace gets value of key in given namespace
func (app *App) GetValueInNamespace(key string, namespace string) (string, error) {
	items, err := app.GetItemsInNamespace(namespace)
	if err != nil {
		return "", err
	}

	return items.Get(key), nil
}

// GetItems gets all the items in default namespace
func (app *App) GetItems() (Items, error) {
	return app.GetItemsInNamespace(defaultNamespace)
}

// GetItemsInNamespace gets all the items in given namespace
func (app *App) GetItemsInNamespace(namespace string) (Items, error) {
	ns, err := app.Client.GetNamespace(namespace, "")
	if err != nil {
		return nil, err
	}

	return ns.Items, nil
}

// GetContent gets the content of given namespace, if the format is properties then will return json string
func (app *App) GetContent(namespace string) (string, error) {
	items, err := app.GetItemsInNamespace(namespace)
	if err != nil {
		return "", err
	}

	if items == nil {
		return "", nil
	}

	if getFormat(namespace) != defaultFormat {
		return items.Get("content"), nil
	}

	return items.String(), nil
}

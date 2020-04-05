package lunar

import (
	"sync"
	"time"
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
	Options                   // inherited options
	ID              string    // app id
	Client          ApolloAPI // the apollo client
	releaseKeyMap   sync.Map  // key: namespace, value: release key
	notificationMap sync.Map  // key: namespace, value: notification id
	watchChan       chan Notification
	errChan         chan error
	stopChan        chan bool
}

// make sure App implements Lunar
var _ Lunar = new(App)

// New creates an application, user must specify the correct app id
func New(appID string, opts ...Option) *App {
	app := &App{
		ID:        appID,
		Options:   NewOptions(opts...),
		watchChan: make(chan Notification),
		errChan:   make(chan error),
		stopChan:  make(chan bool, 1),
	}

	app.UseClient(NewApolloClient(appID, opts...))

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
	items, err := app.GetNamespaceFromApollo(namespace)
	if err != nil {
		return "", err
	}

	return items.Get(key), nil
}

// GetItems gets all the items in default namespace
func (app *App) GetItems() (Items, error) {
	return app.GetNamespaceFromApollo(defaultNamespace)
}

// GetContent gets the content of given namespace, if the format is properties then will return json string
func (app *App) GetContent(namespace string) (string, error) {
	items, err := app.GetNamespaceFromApollo(namespace)
	if err != nil {
		return "", err
	}

	if getFormat(namespace) != defaultFormat {
		return items.Get("content"), nil
	}

	return items.String(), nil
}

// GetItemsInNamespace gets all the items in given namespace.
func (app *App) GetItemsInNamespace(namespace string) (Items, error) {
	return app.GetNamespaceFromApollo(namespace)
}

// GetNamespaceFromApollo gets all the items in given namespace from apollo and refresh cache and map
// This is the most basic method.
func (app *App) GetNamespaceFromApollo(namespace string) (Items, error) {
	namespace = normalizeNamespace(namespace) // trim .properties

	ns, err := app.Client.GetNamespace(namespace, app.getReleaseKey(namespace))
	if err != nil {
		return nil, err
	}

	app.releaseKeyMap.Store(namespace, ns.ReleaseKey)

	// add namespace to notification map with default notification id if not existing,
	// so that it can be watched in long poll
	app.notificationMap.LoadOrStore(namespace, defaultNotificationID)

	// TODO: update cache
	if len(ns.Items) > 0 {

	}

	return ns.Items, nil
}

// gets release key of given namespace
func (app *App) getReleaseKey(namespace string) string {
	if m, ok := app.releaseKeyMap.Load(namespace); ok {
		return m.(string)
	}

	app.releaseKeyMap.Store(namespace, "")

	return ""
}

// Watch watches changes from apollo using long poll
func (app *App) Watch(namespaces ...string) (<-chan Notification, <-chan error) {
	namespaces = refineNamespaces(namespaces)

	// get data from apollo and initialize local namespaces data at the beginning
	for _, namespace := range namespaces {
		app.GetNamespaceFromApollo(namespace)
	}

	// start long poll in goroutine
	go app.startLongPoll()

	return app.watchChan, app.errChan
}

func (app *App) startLongPoll() {
	timer := time.NewTimer(app.LongPollInterval)
	defer timer.Stop()

	for {
		// wait for returns from channel
		select {
		case <-timer.C:
			app.longPoll()
			timer.Reset(app.LongPollInterval)
		case <-app.stopChan:
			app.Logger.Printf("stop watching")
			return
		}
	}
}

// Stop stops watching
func (app *App) Stop() {
	app.stopChan <- true
}

func (app *App) longPoll() error {
	if notifications, err := app.Client.GetNotifications(app.getNotifications()); err == nil {
		// notifications will be empty if no changes
		for _, notification := range notifications {
			// update notification id and then fetch latest data from apollo
			app.notificationMap.Store(notification.Namespace, notification.NotificationID)
			app.GetNamespaceFromApollo(notification.Namespace)

			app.watchChan <- notification
		}
	} else {
		app.Logger.Printf("fail to fetch notifications: %s", err.Error())
		app.errChan <- err

		return err
	}

	return nil
}

func (app *App) getNotifications() Notifications {
	var notifications Notifications

	app.notificationMap.Range(func(key, value interface{}) bool {
		k, _ := key.(string)
		v, _ := value.(int)
		notifications = append(notifications, Notification{
			Namespace:      k,
			NotificationID: v,
		})

		return true
	})

	return notifications
}

# lunar

[![GoDoc](https://godoc.org/github.com/hyperjiang/lunar?status.svg)](https://godoc.org/github.com/hyperjiang/lunar)
[![CI](https://github.com/hyperjiang/lunar/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/hyperjiang/lunar/actions/workflows/ci.yml)
[![](https://goreportcard.com/badge/github.com/hyperjiang/lunar)](https://goreportcard.com/report/github.com/hyperjiang/lunar)
[![codecov](https://codecov.io/gh/hyperjiang/lunar/branch/master/graph/badge.svg)](https://codecov.io/gh/hyperjiang/lunar)
[![Release](https://img.shields.io/github/release/hyperjiang/lunar.svg)](https://github.com/hyperjiang/lunar/releases)

Probably the most elegant ctrip apollo client in golang. This library has no third-party dependency.

Ctrip Apollo: https://github.com/ctripcorp/apollo

Default settings of `lunar`:

- default namespace is `application`
- default server is `localhost:8080`
- default http client timeout is `90s`

*Require golang version >= 1.16 after v0.6.0*

## Usage

```
import "github.com/hyperjiang/lunar"

key := "foo"

app := lunar.New("myAppID", lunar.WithServer("localhost:8080"))

// get value of key in default namespace
app.GetValue(key)

// get value of key in namespace ns
app.GetValueInNamespace(key, "ns")

// get all the items in default namespace
app.GetItems()

// get all the items in namespace ns
app.GetItemsInNamespace("ns")

// get the content of ns namespace, if the format of ns is properties then will return json string
app.GetContent("ns")

// it will fetch items from apollo directly without reading local cache
app.GetNamespaceFromApollo("ns")

// watch changes of given namespaces
watchChan, errChan := app.Watch("ns1", "ns2", ...)

for {
	select {
	case n := <-watchChan:
		fmt.Println(n)
	case <-errChan:
		app.Stop() // stop watcher
		return
	}
}
```

## Logging

`lunar` does not write logs by default, if you want to see logs for debugging, you can replace it with any logger which implements `lunar.Logger` interface.

`lunar` also provide a simple logger `lunar.Printf` which writes to stdout:

```
app := lunar.New("myAppID", lunar.WithServer("localhost:8080"), lunar.WithLogger(lunar.Printf))
```

Or you can use `UseLogger` method:

```
app.UseLogger(lunar.Printf)
```

## Caching

`lunar` use memory cache by default, you can replace it with any cache which implements `lunar.Cache` interface.

`lunar` also provide a file cache `lunar.FileCache` which use files for caching:

```
app := lunar.New("myAppID")

app.UseCache(lunar.NewFileCache("myAppID", "/tmp"))
```

## Enable Access Key

Starting from v1.6.0, apollo supports access key feature, you can use `WithAccessKeySecret` to set the secret:

```
app := lunar.New(
	"myAppID",
	lunar.WithServer("localhost:8080"),
	lunar.WithAccessKeySecret("mySecret"),
)
```

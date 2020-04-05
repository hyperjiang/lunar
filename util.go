package lunar

import (
	"net"
	"path/filepath"
	"strings"
)

func normalizeURL(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	return strings.TrimSuffix(url, "/")
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}

		}
	}

	return ""
}

func getFormat(namespace string) string {
	t := filepath.Ext(namespace)

	if t == "" {
		t = defaultFormat
	}

	return strings.TrimPrefix(t, ".")
}

func normalizeNamespace(namespace string) string {
	return strings.TrimSuffix(namespace, ".properties")
}

func refineNamespaces(namespaces []string) []string {
	type empty struct{}

	namespaces = append(namespaces, defaultNamespace)

	set := make(map[string]empty)
	for _, ns := range namespaces {
		set[normalizeNamespace(ns)] = empty{}
	}

	var result []string
	for k := range set {
		result = append(result, k)
	}

	return result
}

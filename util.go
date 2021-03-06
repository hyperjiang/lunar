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

// GetLocalIP gets local ip
func GetLocalIP() string {
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

func isSupported(format string) bool {
	supportedFormats := []string{".properties", ".xml", ".json", ".yml", ".yaml", ".txt"}

	for _, f := range supportedFormats {
		if format == f {
			return true
		}
	}

	return false
}

// GetFormat gets the format of namespace
func GetFormat(namespace string) string {
	ext := filepath.Ext(namespace)

	if ext == "" || !isSupported(ext) {
		ext = defaultFormat
	}

	return strings.TrimPrefix(ext, ".")
}

// IsProperties checks if the format of namespace is properties
func IsProperties(namespace string) bool {
	return GetFormat(namespace) == defaultFormat
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

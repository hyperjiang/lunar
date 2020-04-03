package lunar

import (
	"net"
	"strings"
)

func normalizeURL(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	return strings.TrimSuffix(url, "/")
}

func splitCommaSeparatedURL(s string) []string {
	var urls []string
	for _, url := range strings.Split(s, ",") {
		urls = append(urls, normalizeURL(strings.TrimSpace(url)))
	}

	return urls
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

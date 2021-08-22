package lunar

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"time"
)

const (
	AuthorizationFormat = "Apollo %s:%s"
	Delimiter           = "\n"
)

func sign(timestamp, pathWithQuery, secret string) string {
	stringToSign := timestamp + Delimiter + pathWithQuery

	h := hmac.New(sha1.New, []byte(secret))
	h.Write([]byte(stringToSign))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func buildHeaders(pathWithQuery, appID, secret string) map[string]string {
	ms := time.Now().UnixNano() / int64(time.Millisecond)
	timestamp := fmt.Sprintf("%d", ms)
	signature := sign(timestamp, pathWithQuery, secret)

	m := make(map[string]string)
	m["Authorization"] = fmt.Sprintf(AuthorizationFormat, appID, signature)
	m["Timestamp"] = timestamp

	return m
}

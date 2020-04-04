package lunar

import "encoding/json"

// Items is items under namespace
type Items map[string]string

// Get gets value of given key
func (items Items) Get(key string) string {
	if v, ok := items[key]; ok {
		return v
	}

	return ""
}

// String converts Items to json string
func (items Items) String() string {
	bytes, _ := json.Marshal(items)
	return string(bytes)
}

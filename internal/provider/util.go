package provider

import "encoding/json"

func ToString(value interface{}) string {
	bytes, _ := json.Marshal(value)
	return string(bytes)
}

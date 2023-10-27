package util

import "encoding/json"

func SerializeToJSONString(v interface{}) (string, error) {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func DeserializeFromJSONString(jsonString string, v interface{}) error {
	return json.Unmarshal([]byte(jsonString), v)
}

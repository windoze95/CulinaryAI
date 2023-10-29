package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
)

func SerializeToJSONString(v interface{}) (string, error) {
	// Check if v is a pointer
	if reflect.ValueOf(v).Kind() != reflect.Ptr {
		return "", errors.New("input must be a pointer")
	}
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func SerializeToJSONStringWithBuffer(v interface{}) (string, error) {
	// Check if v is a pointer
	if reflect.ValueOf(v).Kind() != reflect.Ptr {
		return "", errors.New("input must be a pointer")
	}
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	if err := encoder.Encode(v); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func DeserializeFromJSONString(jsonString string, v interface{}) error {
	// Check if v is a pointer
	if reflect.ValueOf(v).Kind() != reflect.Ptr {
		return errors.New("input must be a pointer")
	}
	return json.Unmarshal([]byte(jsonString), v)
}

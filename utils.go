package main

import (
	"encoding/json"
)

// Convert a given interface to JSON with indentation
func Jsonify(v interface{}) string {
	json_data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(json_data)
}

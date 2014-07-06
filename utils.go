package main

import (
	"encoding/json"
)

func Jsonify(v interface{}) string {
	json_data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(json_data)
}

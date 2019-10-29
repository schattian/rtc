package msh

import "encoding/json"

// Mapable accept any object that can be converted to a key-value pairs version
type Mapable interface {
	ToMap() map[string]interface{}
}

// ToJSON takes a Mapable type and returns the json version of the map
func ToJSON(mapable Mapable) (json.RawMessage, error) {
	mapVersion := mapable.ToMap()
	bytes, err := json.Marshal(mapVersion)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(bytes), nil
}

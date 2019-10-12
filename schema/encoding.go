package schema

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// FromFilename retrieves a schema with the decoded data of the given filename
func FromFilename(filename string) (sch *Schema, err error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	ext := filepath.Ext(filename)
	switch ext {
	case ".json":
		err = json.NewDecoder(bytes.NewReader(body)).Decode(sch)
	case ".yaml":
		err = yaml.NewDecoder(bytes.NewReader(body)).Decode(sch)
	}
	if err != nil {
		return nil, err
	}
	return
}

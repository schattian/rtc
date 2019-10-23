package schema

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// FromFilename retrieves a schema with the decoded data of the given filename
func FromFilename(filename string) (*Schema, error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	sch := &Schema{}
	ext := filepath.Ext(filename)
	switch ext {
	case ".json":
		err = json.NewDecoder(bytes.NewReader(body)).Decode(sch)
	case ".yaml":
		err = yaml.NewDecoder(bytes.NewReader(body)).Decode(sch)
	default:
		err = fmt.Errorf("the extension %v is not allowed", ext)
	}
	if err != nil {
		return nil, err
	}

	err = sch.applyBuiltinValidators()
	if err != nil {
		return nil, err
	}
	return sch, nil
}

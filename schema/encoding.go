package schema

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

func FromFilename(filename string) (sch *Schema, err error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	ext := filepath.Ext(filename)
	switch ext {
	case ".toml":
		_, err = toml.Decode(string(body), sch)
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

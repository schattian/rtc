package schema

import (
	"bytes"
	"encoding/json"
	"path/filepath"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

// FromFilename retrieves a schema with the decoded data of the given filename
func FromFilename(filename string, Fs afero.Fs) (*Schema, error) {
	body, err := afero.ReadFile(Fs, filename)
	if err != nil {
		return nil, err
	}

	sch := &Schema{}
	ext := filepath.Ext(filename)
	switch ext {
	case ".json", ".jsonnet":
		err = json.NewDecoder(bytes.NewReader(body)).Decode(sch)
	case ".yaml":
		err = yaml.NewDecoder(bytes.NewReader(body)).Decode(sch)
	default:
		err = errUnallowedExt
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

package schema

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

// FromFilename retrieves a schema with the decoded data of the given filename
// Notice that the filename can be a remote url
func FromFilename(filename string, Fs afero.Fs) (*Schema, error) {
	var content io.Reader
	var err error

	if strings.HasPrefix(filename, "https") { // Delegate url protocol scheme validation to net/http
		content, err = readRemoteFileByHTTPS(filename)
	} else {
		content, err = readFsFileByName(filename, Fs)
	}
	if err != nil {
		return nil, err
	}

	sch, err := decodeFromReader(content, filepath.Ext(filename))
	if err != nil {
		return nil, err
	}

	err = sch.applyBuiltinValidators()
	if err != nil {
		return nil, err
	}
	return sch, nil
}

func readFsFileByName(filename string, Fs afero.Fs) (io.Reader, error) {
	body, err := afero.ReadFile(Fs, filename)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(body), nil
}

func readRemoteFileByHTTPS(url string) (io.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return resp.Body, nil
}

func decodeFromReader(reader io.Reader, ext string) (sch *Schema, err error) {
	switch ext {
	case ".json", ".jsonnet": // Notice that jsonnet is for already decoded jsonnet files
		err = json.NewDecoder(reader).Decode(sch)
	case ".yaml":
		err = yaml.NewDecoder(reader).Decode(sch)
	default:
		err = errUnallowedExt
	}
	if err != nil {
		return nil, err
	}
	return
}

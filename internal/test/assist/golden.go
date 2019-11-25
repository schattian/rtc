package assist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/google/go-jsonnet"
)

// DecodeJsonnet will load the specified fixture and decode onto the given pointer
// It'll add the boilerplate of testdata/%s.jsonnet
func DecodeJsonnet(name string, pointer interface{}) {
	err := json.Unmarshal(ReadJsonnet(name), pointer)
	if err != nil {
		log.Fatal(err)
	}
}

// ReadJsonnet reads the jsonnet given from the specified file
// It'll add the boilerplate of testdata/%s.jsonnet
func ReadJsonnet(name string) []byte {
	filename := fmt.Sprintf("testdata/%s.jsonnet", name)
	vm := jsonnet.MakeVM()
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error READIN JSONNET FILE: %v ", err)
	}

	json, err := vm.EvaluateSnippet(filename, string(content))
	if err != nil {
		log.Fatalf("Error IN JSONNET EVAL: %v ", err)
	}

	return []byte(json)
}

package assist

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

// DecodeJsonnet will load the specified fixture and decode onto the given pointer
// It'll add the boilerplate of testdata/%s.jsonnet
func DecodeJsonnet(name string, pointer interface{}) {
	err := json.Unmarshal(ReadJsonnet(name), pointer)
	if err != nil {
		log.Fatalf("Error DECODING JSONNET: %v: ", err)
	}
}

// ReadJsonnet reads the jsonnet given from the specified file
// It'll add the boilerplate of testdata/%s.jsonnet
func ReadJsonnet(name string) []byte {
	filename := fmt.Sprintf("testdata/%s.jsonnet", name)
	out, err := exec.Command("jsonnet", filename).Output()
	if err != nil {
		log.Fatalf("Error PARSING JSONNET: %v ", err)
	}
	return out
}

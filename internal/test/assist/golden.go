package assist

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

// DecodeJsonnet will load the specified fixture and decode onto the given pointer
func DecodeJsonnet(name string, pointer interface{}) {
	err := json.Unmarshal(ReadJsonnet(name), pointer)
	if err != nil {
		log.Fatalf("Error DECODING JSONNET: %v: ", err)
	}
}

func ReadJsonnet(name string) []byte {
	filename := fmt.Sprintf("testdata/%s.jsonnet", name)
	out, err := exec.Command("jsonnet", filename).Output()
	if err != nil {
		log.Fatalf("Error PARSING JSONNET: %v ", err)
	}
	return out
}

package assist

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

// DecodeJsonnet will load the specified fixture and decode onto the given pointer
func DecodeJsonnet(name string, pointer interface{}) {
	fileName := fmt.Sprintf("testdata/%s.jsonnet", name)
	out, err := exec.Command("jsonnet", fileName).Output()
	if err != nil {
		log.Fatal(fmt.Errorf("Error parsing jsonnet: %v ", err))
	}

	err = json.Unmarshal(out, pointer)
	if err != nil {
		log.Fatal(fmt.Errorf("Error decoding jsonnet: %v: ", err))
	}
}

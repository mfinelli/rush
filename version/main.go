package version

import (
	"encoding/json"
)

var VERSION version

type version struct {
	Version string `json:"version"`
}

func ParseVersion(pkgjson *[]byte) error {
	if err := json.Unmarshal(*pkgjson, &VERSION); err != nil {
		return err
	}

	return nil
}

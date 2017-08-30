package ipfsapp

import (
	"encoding/json"
	_ "reflect"
)

func struct2json(infoN *InfoNode) ([]byte, error) {
	return json.Marshal(infoN)
}

func json2struct(jsonBytes []byte, infoN *InfoNode) error {
	return json.Unmarshal(jsonBytes, infoN)
}

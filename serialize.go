package ipfsapp

import (
	"encoding/json"
	_ "reflect"
)

func struct2json(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}

func json2struct(jsonstr []byte, obj interface{}) error {
	return nil
}

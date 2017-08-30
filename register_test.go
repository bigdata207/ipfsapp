package ipfsapp

import (
	"encoding/json"
	"fmt"
	"testing"
)

type Human struct {
	name string
	id   string
}

func (h Human) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Name": h.name,
		"ID":   h.id,
	})
}

func (h *Human) UnmarshalJSON(info []byte) error {
	h1 := &struct {
		Name string `json:"Name"`
		ID   string `json:"ID"`
	}{}
	err := json.Unmarshal(info, h1)
	if err != nil {
		return err
	}
	h.name = h1.Name
	h.id = h1.ID
	return nil
}

type Citizen struct {
	h     Human
	stuid string
}

func (s Citizen) MarshalJSON() ([]byte, error) {
	//构造一个字段导出的和Student一样的结构体
	s1 := &struct {
		H     *Human `json:"HumanInfo"`
		StuID string `json:"StuID"`
	}{H: &Human{s.h.name, s.h.id}, StuID: s.stuid}

	return json.Marshal(s1)
}

func (s *Citizen) UnmarshalJSON(info []byte) error {

	s1 := &struct {
		H     Human  `json:"HumanInfo"`
		StuID string `json:"StuID"`
	}{}
	err := json.Unmarshal(info, s1)

	if err != nil {
		return err
	}
	s.h = s1.H
	s.stuid = s1.StuID
	return nil
}

func in2map(m interface{}) {
	fmt.Println(m)
}

func Test(t *testing.T) {
	fmt.Println(structMap)
	registerStruct(Human{})
	fmt.Println(structMap)
	s := Citizen{h: Human{"big", "1234"}, stuid: "223676423745"}
	data, _ := json.Marshal(s)
	in2map(string(data))
}

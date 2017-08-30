package mackzhong

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

type Human struct {
	name string
	id   string
}

func (h Human) String() string {
	return "Name: " + h.name + "\nID: " + h.id + "\n"
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

type Student struct {
	h     *Human
	stuid string
}

func (s Student) MarshalJSON() ([]byte, error) {
	/**
	hi, _ := json.Marshal(s.h)
	return json.Marshal(map[string]interface{}{
		"humaninfo": string(hi),
		"StuID":     s.stuid,
	})
	**/

	//构造一个字段导出的和Student一样的结构体
	s1 := &struct {
		H     *Human `json:"HumanInfo"`
		StuID string `json:"StuID"`
	}{H: &Human{s.h.name, s.h.id}, StuID: s.stuid}
	/**
	s2 := &struct {
		H     *Human `json:"HumanInfo"`
		StuID string `json:"StuID"`
	}{}
	info, _ := json.Marshal(s1)
	json.Unmarshal(info, s2)
	fmt.Println("呵呵")
	fmt.Println(string(info))
	fmt.Println(s2.H)
	fmt.Println("呵呵")
	**/
	return json.Marshal(s1)
}

func (s *Student) UnmarshalJSON(info []byte) error {
	/**
	s1 := &struct {
		HI    string `json:HumanInfo`
		StuID string `json:StuID`
	}{}
	fmt.Println(string(info))
	err := json.Unmarshal(info, s1)
	if err != nil {
		return err
	}
	json.Unmarshal([]byte(s1.HI), s.h)
	s.stuid = s1.StuID
	**/
	s1 := &struct {
		H     *Human `json:"HumanInfo"`
		StuID string `json:"StuID"`
	}{}
	err := json.Unmarshal(info, s1)
	/**
	fmt.Println("求你成功")
	fmt.Println(string(info))
	fmt.Println(s1.H)
	fmt.Println(s1.StuID)
	fmt.Println("求你成功")
	**/
	if err != nil {
		fmt.Println("出事了")
		return err
	}
	s.h = s1.H
	s.stuid = s1.StuID
	return nil
}
func (s Student) String() string {
	return fmt.Sprintf("ID: %v\nName: %v\nStuID: %v\n", s.h.id, s.h.name, s.stuid)
}

type Teacher struct {
	Name string `json:"TName"`
	ID   string `json:"TID"`
}

func (t Teacher) String() string {
	return "Name: " + t.Name + "\nID: " + t.ID + "\n"
}

type Numbers struct {
	ID   int   `json:"id"`
	Nums []int `json:"nums"`
}

type ReHu struct {
	Name string `me:"name"`
	Age  int    `me:"age"`
}

func Test_Json(test *testing.T) {
	fmt.Println("Test JSON!")
	for i := 0; i < 3; i++ {
		fmt.Printf("i = %d\n", i)
	}
	stu := &Student{h: &Human{"zt", "1510"}, stuid: "11510291"}
	fmt.Println(stu)
	h := &Human{"zt", "1510"}
	si, _ := json.Marshal(stu)
	hi, _ := json.Marshal(h)
	fmt.Println(string(si))
	fmt.Println(string(hi))
	t := &Teacher{"zt", "1510"}
	ti, _ := json.Marshal(t)
	t1 := &Teacher{}
	json.Unmarshal(ti, t1)
	fmt.Println(t1)
	h1 := &Human{}
	json.Unmarshal(hi, h1)
	fmt.Println(h1)
	s1 := &Student{}
	err := json.Unmarshal(si, s1)
	if err != nil {
		return
	}
	fmt.Println(s1)
	nn := []int{1, 2, 3}
	ns := &Numbers{1, nn}
	nstr, _ := json.Marshal(ns)
	fmt.Println(string(nstr))
	ns1 := &Numbers{}
	json.Unmarshal(nstr, ns1)
	fmt.Println(ns1.Nums)
	ni, _ := json.Marshal(&Human{name: "zt"})
	fmt.Println(string(ni))
	var hh Human
	json.Unmarshal(ni, &hh)
	fmt.Println(hh.name)
	rh := &ReHu{Name: "zt", Age: 20}
	fmt.Println(rh.Age)
	s := reflect.TypeOf(rh).Elem() //通过反射获取type定义
	//s := &reflect.Type{}
	for i := 0; i < s.NumField(); i++ {
		fmt.Println(s.Field(i).Tag.Get("me")) //将tag输出出来
		fmt.Println(s.Field(i).Name)
		fmt.Println(s.Field(i).Type)
	}
	ss := reflect.ValueOf(rh).Elem()
	ss.FieldByName("Name").SetString("sustc")
	ss.FieldByName("Age").SetInt(18)
	fmt.Println(rh.Age)
	fmt.Println(rh.Name)

	/**
	f, ok := s.FieldByName("Age")
	if ok {
		f.SetInt(18)
	}
	**/
}

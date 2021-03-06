package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strings"

	"os"
	"reflect"
	"strconv"
	"testing"
)

var g2s = map[string]string{"int": "int", "string": "text", "[]byte": "blog"}
var s2g = map[string]string{"int": "int", "text": "string", "blog": "[]byte"}
var DefaultType string = "text"

type Student struct {
	ID     int            `json:"id"`
	Scores map[string]int `json:"scores"`
}

type Class struct {
	TableIndex int                 `json:"tableindex"`
	ClassNum   int                 `json:"classnum"`
	Students   map[string]*Student `json:"students"`
}

func NewClass(TI, CN int, StuJSON string) *Class {
	c := &Class{}
	StuJSON = strings.Replace(StuJSON, "'", "\"", -1)
	tolJSON := fmt.Sprintf(`{"tableindex":%d,"classnum":%d,"students":%s}`, TI, CN, StuJSON)
	fmt.Println(tolJSON)
	json.Unmarshal([]byte(tolJSON), c)
	return c
}

type Demo struct {
	TableIndex int         `json:"tableindex"`
	ClassNum   int         `json:"classnum"`
	ID         int         `json:"id"`
	MM         map[int]int `json:"mm"`
}

func NewDemo(TI, CN, ID int, DemoJSON string) *Demo {
	d := &Demo{}
	DemoJSON = strings.Replace(DemoJSON, "'", "\"", -1)
	tolJSON := fmt.Sprintf(`{"tableindex":%d,"classnum":%d,"id":%d, "mm":%s}`, TI, CN, ID, DemoJSON)
	//fmt.Println(tolJSON)
	json.Unmarshal([]byte(tolJSON), d)
	//fmt.Println(err)
	//fmt.Println(d.MM)
	return d
}

var c1 Class
var d1 Demo

func init() {
	c1.ClassNum = 326
	c1.TableIndex = 326
	d1 = Demo{TableIndex: 326, ClassNum: 326, ID: 326}
	c1.Students = make(map[string]*Student)
	d1.MM = make(map[int]int)
	for i := 0; i < 3; i++ {
		c1.Students[strconv.Itoa(i)] = &Student{ID: i}
		d1.MM[i] = i
		c1.Students[strconv.Itoa(i)].Scores = make(map[string]int)
		k := 0
		for j := 80; j <= 100; j += 10 {
			c1.Students[strconv.Itoa(i)].Scores[strconv.Itoa(100+k)] = j
			k++
		}
	}
}

func CreateTableSQL(oj interface{}, prikey ...interface{}) string {
	var sqlStmt string
	sqlStmt = "create table"
	ot := reflect.TypeOf(oj)
	ov := reflect.ValueOf(oj)
	tableName := ot.Name()
	//fmt.Println(tableName)
	fn := ot.NumField()
	//fmt.Println(fn)
	fields := make(map[string]string)
	for i := 0; i < fn; i++ {
		field := ot.FieldByIndex([]int{i})
		//fmt.Print(field.Name)
		//fmt.Println(": ", field.Type)
		fields[field.Name] = field.Type.String()
	}
	//fmt.Println(oj.(Class).TableIndex)
	val, ok := fields["TableIndex"]
	if ok {
		if strings.Compare(val, "int") == 0 && ov.FieldByName("TableIndex").Int() != 0 {
			//fmt.Println("Yeah")
			//fmt.Println(ov.FieldByName("TableIndex").Int())
			tableName += strconv.Itoa(int(ov.FieldByName("TableIndex").Int()))
		} else {
			if ov.FieldByName("TableIndex").String() != "" {
				tableName += ov.FieldByName("TableIndex").String()
			}
		}
	}
	sqlStmt += (" " + tableName + " (")
	b := len(prikey) > 0
	var pri string = ""
	if b {
		pri += prikey[0].(string)
	}
	for k, v := range fields {
		//fmt.Printf("%s : %s\n", k, v)
		if k != "TableIndex" {
			_, ok := g2s[v]

			if ok {
				sqlStmt += (" " + k + " " + g2s[v])
			} else {
				sqlStmt += (" " + k + " " + DefaultType)
			}
			if pri == k {
				sqlStmt += " primary key"
			}
			sqlStmt += ","
		}
	}
	sqlStmt = sqlStmt[:(len(sqlStmt)-1)] + " )"

	fmt.Println("SQL: ", sqlStmt)
	fmt.Println("TableName: ", tableName)
	return sqlStmt
}

func InsertSQL(oj interface{}) string {
	var sqlStmt string
	sqlStmt = "insert into "
	ot := reflect.TypeOf(oj)
	ov := reflect.ValueOf(oj)
	tableName := ot.Name()
	fmt.Println(tableName)
	fn := ot.NumField()
	//fmt.Println("test")
	fields := make(map[string]string)
	for i := 0; i < fn; i++ {
		field := ot.FieldByIndex([]int{i})
		//fmt.Print(field.Name)
		//fmt.Println(": ", field.Type)
		fields[field.Name] = field.Type.String()
	}
	val, ok := fields["TableIndex"]
	//fmt.Println("test")
	if ok {
		if strings.Compare(val, "int") == 0 && ov.FieldByName("TableIndex").Int() != 0 {
			//fmt.Println("Yeah")
			//fmt.Println(ov.FieldByName("TableIndex").Int())
			tableName += strconv.Itoa(int(ov.FieldByName("TableIndex").Int()))
		} else {
			if ov.FieldByName("TableIndex").String() != "" {
				tableName += ov.FieldByName("TableIndex").String()
			}
		}
	}
	//fmt.Println("test")

	sqlStmt += (" " + tableName + " ( ")
	end := " values ("
	for k, v := range fields {
		//fmt.Printf("%s : %s\n", k, v)
		if k != "TableIndex" {
			var val string
			switch v {
			case "int":
				val = strconv.Itoa(int(ov.FieldByName(k).Int()))
			case "string":
				val = "\"" + ov.FieldByName(k).String() + "\""
			default:
				//val = string(ov.FieldByName(k).Bytes())
				data, _ := json.Marshal(ov.FieldByName(k).Interface())
				val = "\"" + strings.Replace(string(data), "\"", "'", -1) + "\""
			}
			if val != "" {
				sqlStmt += (k + ", ")
				end += (val + ", ")
			}
		}
	}
	//fmt.Println("test")

	sqlStmt = sqlStmt[:(len(sqlStmt)-2)] + " )"
	end = end[:(len(end)-2)] + " )"
	sqlStmt += end
	fmt.Println("SQL: ", sqlStmt)
	return sqlStmt
}
func TestS2M(t *testing.T) {

	info, _ := json.Marshal(&c1)
	//info, _ = bson.Unmarshal(&d1)
	fmt.Println(string(info))
	c2 := &Class{}
	json.Unmarshal(info, c2)
	fmt.Println(c2.ClassNum)
}
func TestM2S(t *testing.T) {
	CreateTableSQL(c1, "ClassNum")
	InsertSQL(c1)
	data, _ := json.Marshal(&c1)
	fmt.Println(string(data))
	cc := &Class{}
	json.Unmarshal(data, cc)
	fmt.Println(cc.ClassNum)
	data, _ = json.Marshal(&d1)
	fmt.Println(string(data))
	dd := &Demo{}
	json.Unmarshal(data, dd)
	fmt.Println(dd.ClassNum)
}

/**
{"tableindex":326,"classnum":326,"students":{"0":{"id":0,"scores":{"100":80,"101":90,"102":100}},"1":{"id":1,"scores":{"100":80,"101":90,"102":100}},"2":{"id":2,"scores":{"100":80,"101":90,"102":100}}}}
**/
func TestOutput2Sqlite3(t *testing.T) {
	fmt.Println("Test Sqlite3 CRUD")
	os.Remove("./foo.db")

	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//defer db.Close()
	sqlStmt := CreateTableSQL(c1, "ClassNum")
	_, err = db.Exec(sqlStmt)
	for i := 0; i < 10; i += 2 {
		c1.ClassNum = i
		sqlStmt = InsertSQL(c1)
		_, err = db.Exec(sqlStmt)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	fmt.Println("Finish")
	fmt.Println("****************Select*******************8")
	rows, err := db.Query("select ClassNum,Students from Class326")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var ClassNum int
		var Students string
		err = rows.Scan(&ClassNum, &Students)
		cc := NewClass(326, ClassNum, Students)
		if err != nil {
			fmt.Println(err)
			return
		}
		data, _ := json.Marshal(cc)
		fmt.Println(string(data))
	}
	fmt.Println("Try Demo")
	sqlStmt = CreateTableSQL(d1, "ClassNum")
	_, err = db.Exec(sqlStmt)
	for i := 0; i < 10; i += 2 {
		d1.ClassNum = i
		sqlStmt = InsertSQL(d1)
		_, err = db.Exec(sqlStmt)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	fmt.Println("Finish")
	fmt.Println("****************Select*******************8")
	rows.Close()
	rows, err = db.Query("select * from Demo326")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var ClassNum int
		var ID int
		var MM string
		err = rows.Scan(&ClassNum, &ID, &MM)
		dd := NewDemo(326, ClassNum, ID, MM)
		if err != nil {
			fmt.Println(err)
			return
		}
		data, _ := json.Marshal(dd)
		fmt.Println(string(data))
	}
}

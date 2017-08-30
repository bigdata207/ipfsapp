package ipfsapp

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"testing"
)

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

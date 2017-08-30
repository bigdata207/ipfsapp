package ipfsapp

import (
	"encoding/json"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	_ "github.com/syndtr/goleveldb/leveldb/opt"
	"testing"
)

type Sustc struct {
	Name string
	Age  int
}

func TestLevelDB(t *testing.T) {
	fmt.Println("Test LevelDB")
	ss := Sustc{Name: "t", Age: 8}
	info, _ := json.Marshal(ss)
	fmt.Println(string(info))
	//o := &opt.Options{}
	db, err := leveldb.OpenFile("db", nil)
	if err != nil {
		return
	}
	defer db.Close()
	//db.Put([]byte(ss.Name), info, nil)
	db.Delete([]byte("t"), nil)
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		key := iter.Key()
		value := iter.Value()
		fmt.Printf("%s : %s\n", key, value)
	}
	me, _ := db.Get([]byte("hsg"), nil)
	fmt.Printf("%s\n", me)
	iter.Release()
	err = iter.Error()
}

package ipfsapp

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	_ "github.com/syndtr/goleveldb/leveldb/opt"
)

type userConn struct {
	name   string
	pubKey string
	priKey string
	dbName string
	db     leveldb.DB
}

func (user userConn) Get(key string) string {
	return fmt.Sprintf("Get %v", key)
}
func (user userConn) Update(key string) string {
	return fmt.Sprintf("Update %v", key)
}
func (user *userConn) Add(key, value string) error {
	return nil
}
func (user *userConn) Delete(key string) error {
	return nil
}

type adminConn struct {
	name   string
	pubKey string
	priKey string
	dbName string
	db     leveldb.DB
}

func (a adminConn) getIndexTree(key string) string {
	return "get " + key
}

func (a *adminConn) updateIndexTree(string) error {
	return nil
}

func (a *adminConn) add(key, tree string) error {
	return nil
}
func (a *adminConn) delete(key string) error {
	return nil
}

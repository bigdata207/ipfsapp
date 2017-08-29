package main

import "fmt"

type userConn struct{}

func Get(key string) string {
	return fmt.Sprintf("hello %v", key)
}

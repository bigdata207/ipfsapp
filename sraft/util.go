package sraft

import "log"

//Debug Debug标识符
const Debug = 1

//DPrintf 日志输出debug信息
func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug > 0 {
		log.Printf(format, a...)
	}
	return
}

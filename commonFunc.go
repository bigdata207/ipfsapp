package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func readLine() ([]byte, error) {
	reader := bufio.NewReader(os.Stdin)
	data, b, err := reader.ReadLine()
	fmt.Println(b)
	return data, err
}
func relaPath2AbsoPath(relaPath string) string {
	curDir, _ := os.Getwd()
	if relaPath[0] == '.' {
		return curDir + "/" + relaPath[1:]
	} else if relaPath[0] != '/' {
		return curDir + "/" + relaPath
	}
	return relaPath
}
func search(in IndexNode, key string) []Path {
	paths := make([]Path, 0)
	if strings.Contains(in.kv.k, key) {
		paths = append(paths, in.kv.v)
	}

	if in.isDir() {
		for _, v := range in.childBlocks {
			tmp := search(v, key)
			if len(tmp) == 0 {
				return paths
			}
			paths = append(paths, tmp...)
		}
	}
	return paths
}

func dfs(kv kvNode) map[string]IndexNode {
	dirPath := kv.v
	ins := make(map[string]IndexNode)
	childs, err := ioutil.ReadDir(dirPath.toString())
	if err == nil {
		for _, v := range childs {
			tmp := deepth(dirPath.Add(v.Name()))
			diff := tmp.Sub(dirPath)
			if len(diff) > 0 {
				ins[diff] = NewIndexNode(DirNode{}, kvNode{k: diff, v: tmp})
			}
		}
	}
	return ins
}

func deepth(nodePath Path, e ...interface{}) Path {
	childs, err := ioutil.ReadDir(nodePath.toString())
	if err == nil {
		if len(childs) == 1 && childs[0].IsDir() {
			childPath, err := nodePath.Add(childs[0].Name())
			if err == nil {
				return deepth(childPath)
			}
		}
	}
	return ""
}

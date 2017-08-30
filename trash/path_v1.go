package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	_ "path/filepath"
	"strings"
	"time"
)

var rootDir = flag.String("dir", "/mnt/extra/MackZ/go", "The root directure to generate Dircture Hash Tree")

var isDir = map[bool]string{false: "File", true: "Directure"}

//FileNode is a Node struct of a file or directure
type FileNode struct {
	Pre      *FileNode
	Next     map[string]*FileNode
	Name     string
	FullPath string
	Version  string
	Level    int
	Type     string
	IPFS     map[string]string
}

//NewFileNode Return a FileNodeointer
func NewFileNode(Path string, b bool) (fn *FileNode) {
	fn = &FileNode{FullPath: Path}
	pathSeq := strings.Split(Path, string(os.PathSeparator))
	fn.Level = len(pathSeq) - 1
	fn.Type = isDir[b]
	fn.Name = pathSeq[len(pathSeq)-1]
	fn.Next = make(map[string]*FileNode)
	fn.IPFS = make(map[string]string)
	return
}

//解析成父节点，当前节点，后续节点
// MarshalJson is used to implement the json rule
func (fn FileNode) MarshalJSON() ([]byte, error) {
	//fmt.Println(len(fn.Next))
	childPath := make([]string, 0, len(fn.Next))
	for _, n := range fn.Next {
		childPath = append(childPath, n.Name)
	}
	return json.Marshal(map[string]interface{}{
		"Name":     fn.Name,
		"FullPath": fn.FullPath,
		"Version":  fn.Version,
		"Level":    fn.Level,
		"Type":     fn.Type,
		"Next":     childPath,
		"IPFS":     fn.IPFS,
	})
}

func GenerateMD5(data []byte) (md5Key string) {
	md5Ctx := md5.New()
	md5Ctx.Write(data)
	cipherStr := md5Ctx.Sum(nil)
	//fmt.Println(cipherStr)
	//fmt.Println(hex.EncodeToString(cipherStr))
	md5Key = hex.EncodeToString(cipherStr)
	return
}

//HashTree use to find the FileNode fastly
type HashTree struct {
	htmap map[string]*FileNode
	head  *FileNode
}

func (ht HashTree) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"head":  ht.head.FullPath,
		"htmap": ht.htmap,
	})
}

func (ht *HashTree) put(dirName string, fd *FileNode) {
	md5Key := GenerateMD5([]byte(dirName))
	ht.htmap[md5Key] = fd
}

func (ht HashTree) Get(dirName string) (*FileNode, bool) {
	md5Key := GenerateMD5([]byte(dirName))
	fd, ok := ht.htmap[md5Key]
	if ok {
		return fd, ok
	}
	return nil, ok
}

func NewHashTree() *HashTree {
	ht := &HashTree{}
	ht.htmap = make(map[string]*FileNode)
	return ht
}

//var ht HashTree

//BuildTree is used to get the directure tree of the special directure path
func (ht *HashTree) BuildTree(pre *FileNode, dirPath string, b bool, suffix ...string) (root *FileNode, err error) {
	fns := NewFileNode(dirPath, b)
	fns.Pre = pre
	hasSuffix := (len(suffix) > 0)
	//fmt.Println(hasSuffix)
	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	PthSep := string(os.PathSeparator)
	var sf string
	if hasSuffix {
		sf = strings.ToUpper(suffix[0]) //忽略后缀匹配的大小写
	}

	for _, fi := range dir {
		fileFullPath := dirPath + PthSep + fi.Name()
		if fi.IsDir() { // 忽略目录
			nextNode, err := ht.BuildTree(fns, fileFullPath, true, suffix...)
			if err != nil {
				log.Fatalln(err)
			}
			if nextNode != nil {
				//fns.Next = append(fns.Next, nextNode)
				fns.Next[GenerateMD5([]byte(fileFullPath))] = nextNode
			}
		} else {
			if hasSuffix { //匹配文件}
				if strings.HasSuffix(strings.ToUpper(fi.Name()), sf) {
					nextNode := NewFileNode(fileFullPath, false)
					nextNode.Pre = fns
					//fns.Next = append(fns.Next, nextNode)
					fns.Next[GenerateMD5([]byte(fileFullPath))] = nextNode
					//fmt.Println(fi.Name())
					ht.put(fileFullPath, nextNode)
				}
			} else {
				nextNode := NewFileNode(fileFullPath, false)
				nextNode.Pre = fns
				fns.Next[GenerateMD5([]byte(fileFullPath))] = nextNode
				ht.put(fileFullPath, nextNode)
			}
		}
	}
	ht.put(fns.FullPath, fns)
	//fmt.Println("Success")
	return fns, nil
}

func (ht *HashTree) freeNode(fn *FileNode) {
	for _, n := range fn.Next {
		ht.freeNode(n)
	}
	delete(ht.htmap, GenerateMD5([]byte(fn.Name)))
}

func (ht *HashTree) DeleteFileNode(nodePath string) {
	md5Key := GenerateMD5([]byte(nodePath))
	//fmt.Printf("%v : %v\n", nodePath, md5Key)
	//fmt.Println("Test")
	preMD5 := GenerateMD5([]byte(ht.htmap[md5Key].Pre.FullPath))
	//fmt.Println("Test")
	delete(ht.htmap[preMD5].Next, md5Key)
	//fmt.Println("Test")
	ht.freeNode(ht.htmap[md5Key])
}

func (ht *HashTree) AddFileNode(nodePath string, b bool) {
	pathSeq := strings.Split(nodePath, string(os.PathSeparator))
	preNodePath := nodePath[:(len(nodePath) - len(pathSeq[len(pathSeq)-1]) - 1)]
	//fmt.Println(preNodePath)
	md5Key := GenerateMD5([]byte(nodePath))
	preKey := GenerateMD5([]byte(preNodePath))
	newNode, err := ht.BuildTree(ht.htmap[preKey], nodePath, b)
	if err == nil {
		ht.htmap[preKey].Next[md5Key] = newNode
	}
}

func (ht *HashTree) MoveFileNode(nodePath, newPreDir string) bool {
	pathSeq := strings.Split(nodePath, string(os.PathSeparator))
	fileName := pathSeq[len(pathSeq)-1]
	md5Key := GenerateMD5([]byte(nodePath))
	preMD5 := GenerateMD5([]byte(ht.htmap[md5Key].Pre.FullPath))
	newPreKey := GenerateMD5([]byte(newPreDir))
	//fmt.Println("Test")
	delete(ht.htmap[preMD5].Next, md5Key)
	ht.htmap[md5Key].Pre = ht.htmap[newPreKey]
	//fmt.Println("Test")

	ht.htmap[newPreKey].Next[md5Key] = ht.htmap[md5Key]
	//fmt.Println("Test")

	ht.htmap[md5Key].FullPath = newPreDir + string(os.PathSeparator) + fileName
	//fmt.Println("Test")

	_, ok := ht.htmap[newPreKey].Next[md5Key]

	return ok
}

func (ht *HashTree) String() string {
	jsonstr, _ := json.Marshal(ht)
	return string(jsonstr)
}

func (ht HashTree) FuzzySearch(keyWord string) (paths []string) {
	for _, v := range ht.htmap {
		if strings.Contains(v.FullPath, keyWord) {
			paths = append(paths, v.FullPath)
		}
	}
	return
}

type User struct {
	Git            map[string]*HashTree
	Versions       []string
	CurrentVersion string
}

func NewUser() *User {
	u := &User{}
	u.Git = make(map[string]*HashTree)
	return u
}
func (u *User) NewVersion(rootDir, version string) {
	u.Git[version] = NewHashTree()
	u.Versions = append(u.Versions, version)
	//fmt.Println(len(u.Git))
	_, err := ioutil.ReadDir(rootDir)
	if err == nil {
		u.Git[version].head, _ = u.Git[version].BuildTree(nil, rootDir, true)
	} else {
		u.Git[version].head, _ = u.Git[version].BuildTree(nil, rootDir, false)
	}
}

func (u *User) SwitchVersion(version string) {
	if _, ok := u.Git[version]; ok {
		u.CurrentVersion = version
		return
	}
	fmt.Println("No this version!")
}
func (u *User) DeleteVersion(version string) {
	if _, ok := u.Git[version]; ok {
		delete(u.Git, version)
		return
	}
	fmt.Println("No this version!")
}

//*rootDir = "/mnt/extra/MackZ/go"
var demoDir string
var demoFile string

func (u User) Test(version ...string) {
	var ok bool
	var ht *HashTree
	var v string
	if len(version) > 0 {
		v = version[0]
		ht, ok = u.Git[version[0]]
	} else {
		v = u.CurrentVersion
		ht, ok = u.Git[u.CurrentVersion]
	}
	//fmt.Println(ok)
	if ok {
		show := false
		t2 := time.Now()
		_, err := ioutil.ReadDir(*rootDir)
		if err != nil {
			log.Fatalln(err)
		}
		root, err := ht.BuildTree(nil, *rootDir, true)
		ht.head = root
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("Root: ", root.FullPath)
		node, _ := ht.Get(root.FullPath)
		fmt.Printf("Befor delete(%v): %d\n", root.FullPath, len(node.Next))
		if show {
			for i, v := range node.Next {
				fmt.Printf("%v : %v\n", i, v.FullPath)
			}
		}

		ht.DeleteFileNode(demoDir)
		node, _ = ht.Get(root.FullPath)
		fmt.Printf("After delete(%v): %d\n", root.FullPath, len(node.Next))
		if show {
			for i, v := range node.Next {
				fmt.Printf("%v : %v\n", i, v.FullPath)
			}
		}

		ht.AddFileNode(demoDir, true)
		node, _ = ht.Get(root.FullPath)
		fmt.Printf("After add(%v): %d\n", root.FullPath, len(node.Next))
		if show {
			for i, v := range node.Next {
				fmt.Printf("%v : %v\n", i, v.FullPath)
			}
		}

		node, _ = ht.Get(demoDir)
		fmt.Printf("Before move(%v): %d\n", demoDir, len(node.Next))
		if show {
			for i, v := range node.Next {
				fmt.Printf("%v : %v\n", i, v.FullPath)
			}
		}
		ok := ht.MoveFileNode(demoFile, root.FullPath)
		node, _ = ht.Get(root.FullPath)
		node, _ = ht.Get(root.FullPath)
		fmt.Printf("After move(%v): %d\n", root.FullPath, len(node.Next))
		if show {
			for i, v := range node.Next {
				fmt.Printf("%v : %v\n", i, v.FullPath)
			}
		}

		node, _ = ht.Get(demoDir)
		fmt.Printf("After move(%v): %d\n", demoDir, len(node.Next))
		if show {
			for i, v := range node.Next {
				fmt.Printf("%v : %v\n", i, v.FullPath)
			}
		}
		fmt.Println(ok)
		t1 := time.Now()
		paths := ht.FuzzySearch("chacha_test.go")
		fmt.Println("Running Time for Searching: ", time.Since(t1))
		for i, v := range paths {
			fmt.Printf("%d : %v\n", i, v)
		}
		RunningTime := time.Since(t2)
		fmt.Println("Running Time from Begin: ", RunningTime)
	} else {
		fmt.Println("No this version: ", v)
	}
}

func (ht *HashTree) InitFromDatabase() bool {
	b := false

	return b
}
func (u *User) SaveToDatabase() bool {
	b := false
	return b
}

func main() {
	flag.Parse()
	demoDir = *rootDir + "/bin"
	demoFile = demoDir + "/geth"
	version := "sustc"
	u := NewUser()

	u.NewVersion(*rootDir, version)
	u.SwitchVersion(version)
	u.Test()
	fmt.Println(u.Git[version].InitFromDatabase())

	//fmt.Println(len(u.Git[version].head.Next))
	info, _ := json.Marshal(u.Git[version].head)
	fmt.Println(string(info))
	info, _ = json.Marshal(u.Git[version].htmap)
	ioutil.WriteFile("DirTree.txt", info, 0644)
	//fmt.Println(string(info))
	u.DeleteVersion(version)
	u.Test()
}

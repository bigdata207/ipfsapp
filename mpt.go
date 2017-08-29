package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	_ "reflect"
	_ "strconv"
	"strings"
)

const (
	user       = "zt"
	password   = "1234"
	seq        = "!@#$%^&*()_+"
	keyRoot    = "./"
	priKeyFile = "private.pem"
	pubKeyFile = "public.pem"
)

var salt []byte
var privateKey []byte
var publicKey []byte

func calSalt() {
	files, err := ioutil.ReadDir(keyRoot)
	if err == nil {
		b := 0
		for _, v := range files {
			if !v.IsDir() && (strings.Compare(v.Name(), priKeyFile) == 0 || strings.Compare(v.Name(), pubKeyFile) == 0) {
				b++
			}
			fmt.Println(v.Name())
		}
		if b != 2 {
			fmt.Println("Can't find key file, regenerate!")
			genRSAKey()
		}
		privateKey, _ = ioutil.ReadFile(keyRoot + priKeyFile)

		publicKey, _ = ioutil.ReadFile(keyRoot + pubKeyFile)
	}
	origstr := user + seq + password
	if len([]byte(origstr)) > 128 {
		fmt.Println("seq is too long")
		return
	}
	salt, err = RsaEncrypt([]byte(origstr))
}

type Path string

func NewPath(p string) (Path, error) {
	if p[0] != '/' {
		//fmt.Println("path should begin with '/'")
		return "", errors.New("path should begin with '/'")
	}
	return Path(p), nil
}

func (p Path) Add(end string) (Path, error) {
	ps := p.toString()
	if end[0] == '/' {
		return p, errors.New("sub path shouldn't begin with '/'")
	}
	tmp := ""
	//fmt.Println(p[len(p)-1] == '/')
	if ps[len(p)-1] == '/' {
		tmp += (ps + end)
	} else {
		tmp += (ps + "/" + end)
	}
	return Path(tmp), nil
}

func (p Path) Back(end string) (Path, error) {
	ps := p.toString()
	//if end[0] == '/' {
	//	return p, errors.New("sub path shouldn't begin with '/'")
	//}
	ok := strings.Contains(ps, end)
	if ok && len(ps) > len(end) {
		if ps[(len(p)-len(end))-1] == '/' {
			return Path(ps[:(len(p) - len(end))]), nil
		}
	}
	return p, errors.New("Miss match")
}

func (p Path) Sub(ap interface{}) string {
	switch ap.(type) {
	case string:
		b := strings.Compare(ap.(string), p.toString()[:len(ap.(string))]) == 0
		if b {
			return p.toString()[len(ap.(string)):]
		}
		return ""
	case Path:
		tmp := ap.(Path).toString()
		b := strings.Compare(tmp, p.toString()[:len(tmp)]) == 0
		if b {
			return p.toString()[len(ap.(string)):]
		}
		return ""
	default:
		return ""
	}
}
func (p Path) KeyHash() string {
	data := md5.Sum(append([]byte(p), salt...))
	return fmt.Sprintf("%x", data)
}

func (p Path) toString() string {
	return string(p)
}

type parent struct {
	Hash         string
	RelativePath Path
}

type piece struct {
	Index         int
	PieceName     string
	PieceIPFSHash string
	Log           []byte
}

type InfoNode struct {
	FileName       string
	UploadRoot     parent //Reord the hash and the relative path of the top folder at the latest upload
	AbsolutePath   Path
	IPFSHash       string
	Pieces         piece    //record the ipfshash of each pieces
	PathAndContent string   //Hash(ipfshash +hash(path)), if path or content changed, this value will change
	ChildBlocks    []string //Child Files and Folders
	CryptoKey      string   //use to crypto the pieces
	//UploadType     string   //Client or Web
	FileType string //Directure Or File
	MD5      string
}

type kvNode struct {
	k string
	v Path
}

func (kv kvNode) Get(key string) (Path, error) {
	if strings.Compare(kv.k, key) == 0 {
		return kv.v, nil
	}
	return "", errors.New("No such node")
}

type childNodes struct {
}
type NodeType interface {
	toString() string
}
type DirNode struct{}

func (DirNode) toString() string {
	return "Directure"
}

type FileNode struct{}

func (FileNode) toString() string {
	return "File"
}

type IndexNode struct {
	nType       NodeType
	kv          kvNode
	childBlocks map[string]IndexNode
}

func NewIndexNode(nt NodeType, kv kvNode) IndexNode {
	return IndexNode{nType: nt, kv: kv}
}

func (in IndexNode) isDir() bool {
	return strings.Compare(in.nType.toString(), "Directure") == 0
}

func (in IndexNode) Type() string {
	return in.nType.toString()
}

func (in IndexNode) Download() error {
	return nil
}

package ipfsapp

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	_ "reflect"
	_ "strconv"
	"strings"
)

/**常量参数
* user: 用户名
* password: 用户密码
* seq:用户名和密码中间的分割串
* ketRoot: rsa key文件存放目录
* priKeyFile:rsa私钥文件名
* pubKeyFile: rsa公钥文件名
**/

const (
	username   = "bigdata"
	password   = "1234"
	seq        = "!@#$%^&*()_+"
	keyRoot    = "./"
	priKeyFile = "private.pem"
	pubKeyFile = "public.pem"
)

//将账户密码生成特定串经rsa加密后加盐hash生成leveldb　key
var salt []byte

//rsa私钥
var privateKey []byte

//rsa公钥
var publicKey []byte

//初始化公钥和私钥
func init() {
	files, err := ioutil.ReadDir(keyRoot)
	if err == nil {
		b := 0
		for _, v := range files {
			if !v.IsDir() && (strings.Compare(v.Name(), priKeyFile) == 0 || strings.Compare(v.Name(), pubKeyFile) == 0) {
				b++
			}
			//fmt.Println(v.Name())
		}
		if b != 2 {
			fmt.Println("Can't find key file, regenerate!")
			genRSAKey()
		}
		privateKey, _ = ioutil.ReadFile(keyRoot + priKeyFile)

		publicKey, _ = ioutil.ReadFile(keyRoot + pubKeyFile)
	}
}

//calSalt 计算salt,并且会在不存在rsa key时生成
func calSalt() {

	origstr := username + seq + password
	if len([]byte(origstr)) > 128 {
		fmt.Println("seq is too long")
		return
	}
	salt, err = RsaEncrypt([]byte(origstr), publicKey)
}

//Path 基于string的Path结构
type Path string

//NewPath 根据给定的路径字符串返回一个Path对象
func NewPath(p string) (Path, error) {
	if p[0] != '/' {
		//fmt.Println("path should begin with '/'")
		return "", errors.New("path should begin with '/'")
	}
	return Path(p), nil
}

//Add 从当前前进几层目录
func (p Path) Add(end string) (Path, error) {
	ps := p.String()
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

//Back 从当前回退几层目录
func (p Path) Back(end string) (Path, error) {
	ps := p.String()
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

//Sub 一个path对象减去另一个path对象或者字符串，要求另一个path对象是该对象的父路径
func (p Path) Sub(ap interface{}) string {
	switch ap.(type) {
	case string:
		b := strings.Compare(ap.(string), p.String()[:len(ap.(string))]) == 0
		if b {
			return p.String()[len(ap.(string)):]
		}
		return ""
	case Path:
		tmp := ap.(Path).String()
		b := strings.Compare(tmp, p.String()[:len(tmp)]) == 0
		if b {
			return p.String()[len(ap.(string)):]
		}
		return ""
	default:
		return ""
	}
}

//KeyHash 生成LevelDB中查询文件信息所用的Key
func (p Path) KeyHash() string {
	data := md5.Sum(append([]byte(p), salt...))
	return fmt.Sprintf("%x", data)
}

//String path对象返回其string字符串
func (p Path) String() string {
	return string(p)
}

type parent struct {
	Hash         string
	RelativePath Path
}

//NewParent 返回一个上传记录对象
func NewParent(h string, r Path) parent {
	return parent{Hash: h, RelativePath: r}
}

type piece struct {
	Index     int
	PieceName string
	PieceIPFS string
	Log       []byte
}

//NewPiece 返回一个分片对象
func NewPiece(i int, pn, pi string, l []byte) piece {
	return piece{Index: i, PieceName: pn, PieceIPFS: pi, Log: l}
}

//InfoNode 文件信息节点,序列化后存入leveldb,从leveldb中读取的数据也会反序列化成InfoNode对象
type InfoNode struct {
	FileName       string   //文件[夹]名
	UploadRoot     parent   //记录文件上传时的根文件夹的hash以及与其相对路径
	AbsolutePath   Path     //文件[夹]在目录树中的绝对路径
	IPFSHash       string   //文件直接上传的IPFS链接
	Pieces         []piece  //记录RS分块信息与IPFSHash二选一
	PathAndContent string   //Hash(ipfshash +hash(path)),路径或者内容改变这个值都会改变
	ChildBlocks    []string //子文件[夹]
	CryptoKey      string   //文件AES加密密钥
	//UploadType     string   //Client or Web
	FileType string //文件类型
	MD5      string //记录文件md5校验值
}

//InfoNode2Json 将InfoNode对象转成json字符串
func InfoNode2Json(infoN *InfoNode) ([]byte, error) {
	return json.Marshal(infoN)
}

//NewInfoNodeFromJson 从json对象中恢复一个InfoNode对象
func NewInfoNodeFromJson(jsonBytes []byte) (*InfoNode, error) {
	in := &InfoNode{}
	err := json.Unmarshal(jsonBytes, in)
	return in, err
}

//kvNode　索引树中的key:绝对路径
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

//NodeType 文件类型接口,实现该接口的String()返回文件类型字符串File or Directure
type NodeType interface {
	String() string
}

//DirNode 实现NodeType接口，返回一个"Directure"
type DirNode struct{}

func (DirNode) String() string {
	return "Directure"
}

//FileNode 实现NodeType接口，返回一个"File"
type FileNode struct{}

func (FileNode) String() string {
	return "File"
}

//IndexNode 索引结构体定义
type IndexNode struct {
	nType       NodeType
	kv          kvNode
	childBlocks map[string]IndexNode
}

//NewIndexNode 根据传入的NodeType和kvNode新建一个IndexNode并返回
func NewIndexNode(nt NodeType, kv kvNode) IndexNode {
	return IndexNode{nType: nt, kv: kv}
}

//IsDir 返回节点是否是文件夹节点
func (in IndexNode) IsDir() bool {
	return strings.Compare(in.nType.String(), "Directure") == 0
}

//Type 返回节点类型字符串
func (in IndexNode) Type() string {
	return in.nType.String()
}

func (in IndexNode) Download() error {
	return nil
}

package ipfsapp

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	bhost "github.com/libp2p/go-libp2p/p2p/host/basic"
	pstore "gx/ipfs/QmPgDWmTmuzvP7QE5zwo1TmjbJme9pmZHNujB2453jkCTr/go-libp2p-peerstore"
	gologging "gx/ipfs/QmQvJiADDe7JR4m968MwXobTCCzUqQkP87aRHe29MEBGHV/go-logging"
	msmux "gx/ipfs/QmRVYfZ7tWNHPBzWiG6KWGzvT2hcGems8srihsQE29x1U5/go-smux-multistream"
	golog "gx/ipfs/QmSpJByNKFX1sCsHBEp3R73FL4NF6FnQTEGyNAXHm2GS52/go-log"
	swarm "gx/ipfs/QmW8QNRUf1nqeRxZdhGg5DVcxHMwxuuNt7AWfmDi2a8JE2/go-libp2p-swarm"
	ma "gx/ipfs/QmXY77cVe7rVRQXZZQRioukUM7aRW3BTcAgJe12MCtb3Ji/go-multiaddr"
	peer "gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
	crypto "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	net "gx/ipfs/QmahYsGWry85Y7WUe2SX5G4JkH2zifEQAUtJVLZ24aC9DF/go-libp2p-net"
	yamux "gx/ipfs/Qmbn7RYyWzBVXiUp9jZ1dA4VADHy9DtS7iZLwfhEUQvm3U/go-smux-yamux"
	host "gx/ipfs/Qmeqtv7ASvepCpUDxysK2oP1urtEr5eMNMxbhTJYJziAGo/go-libp2p-host"
	"io"
	"io/ioutil"
	"log"
	mrand "math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// makeBasicHost creates a LibP2P host with a random peer ID listening on the
// given multiaddress. It will use secio if secio is true.
func makeBasicHosts(listenPort int, secio bool, randseed int64) (map[string]host.Host, error) {

	// If the seed is zero, use real cryptographic randomness. Otherwise, use a
	// deterministic randomness source to make generated keys stay the same
	// across multiple runs
	var r io.Reader
	if randseed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(randseed))
	}

	// Generate a key pair for this host. We will use it at least
	// to obtain a valid host ID.
	priv, pub, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return nil, err
	}

	// Obtain Peer ID from public key
	pid, err := peer.IDFromPublicKey(pub)
	if err != nil {
		return nil, err
	}
	ips := getInternalIP()
	//ips = append(ips, "116.7.234.243")
	bhs := make(map[string]host.Host)
	for _, ip := range ips {

		// Create a multiaddress
		addr, err := ma.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", ip, listenPort))

		if err != nil {
			return nil, err
		}

		// Create a peerstore
		ps := pstore.NewPeerstore()

		// If using secio, we add the keys to the peerstore
		// for this peer ID.
		if secio {
			ps.AddPrivKey(pid, priv)
			ps.AddPubKey(pid, pub)
		}

		// Set up stream multiplexer
		tpt := msmux.NewBlankTransport()
		tpt.AddTransport("/yamux/1.0.0", yamux.DefaultTransport)

		// Create swarm (implements libP2P Network)
		swrm, err := swarm.NewSwarmWithProtector(
			context.Background(),
			[]ma.Multiaddr{addr},
			pid,
			ps,
			nil,
			tpt,
			nil,
		)
		if err != nil {
			return nil, err
		}

		netw := (*swarm.Network)(swrm)

		basicHost := bhost.New(netw)
		bhs[ip] = basicHost
		// Build host multiaddress
		hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", basicHost.ID().Pretty()))

		// Now we can build a full multiaddress to reach this host
		// by encapsulating both addresses:
		fullAddr := addr.Encapsulate(hostAddr)
		log.Printf("I am %s\n", fullAddr)
		if secio {
			log.Printf("Now run \"./echo -l %d -d %s -secio\" on a different terminal\n", listenPort+1, fullAddr)
		} else {
			log.Printf("Now run \"./echo -l %d -d %s\" on a different terminal\n", listenPort+1, fullAddr)
		}
	}
	return bhs, nil
}
func makeBasicHost(listenPort int, secio bool, randseed int64) (host.Host, error) {

	// If the seed is zero, use real cryptographic randomness. Otherwise, use a
	// deterministic randomness source to make generated keys stay the same
	// across multiple runs
	var r io.Reader
	if randseed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(randseed))
	}

	// Generate a key pair for this host. We will use it at least
	// to obtain a valid host ID.
	priv, pub, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return nil, err
	}

	// Obtain Peer ID from public key
	pid, err := peer.IDFromPublicKey(pub)
	if err != nil {
		return nil, err
	}

	// Create a multiaddress
	addr, err := ma.NewMultiaddr(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort))

	if err != nil {
		return nil, err
	}

	// Create a peerstore
	ps := pstore.NewPeerstore()

	// If using secio, we add the keys to the peerstore
	// for this peer ID.
	if secio {
		ps.AddPrivKey(pid, priv)
		ps.AddPubKey(pid, pub)
	}

	// Set up stream multiplexer
	tpt := msmux.NewBlankTransport()
	tpt.AddTransport("/yamux/1.0.0", yamux.DefaultTransport)

	// Create swarm (implements libP2P Network)
	swrm, err := swarm.NewSwarmWithProtector(
		context.Background(),
		[]ma.Multiaddr{addr},
		pid,
		ps,
		nil,
		tpt,
		nil,
	)
	if err != nil {
		return nil, err
	}

	netw := (*swarm.Network)(swrm)

	basicHost := bhost.New(netw)

	// Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", basicHost.ID().Pretty()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	fullAddr := addr.Encapsulate(hostAddr)
	log.Printf("I am %s\n", fullAddr)
	if secio {
		log.Printf("Now run \"./echo -l %d -d %s -secio\" on a different terminal\n", listenPort+1, fullAddr)
	} else {
		log.Printf("Now run \"./echo -l %d -d %s\" on a different terminal\n", listenPort+1, fullAddr)
	}

	return basicHost, nil
}

var mb int = 1024 * 1024

type p2prequest struct {
	Op        string `json:"op"`
	DataSize  int    `json:"datasize"`
	HasPubKey bool   `json:"haspubkey"`
	PubKey    string `json:"pubkey"`
}

type p2preply struct {
	Status  bool   `json:"status"`
	GetSize int    `json:"getsize"`
	Op      string `json:"op"`
	MD5     string `json:"md5"`
}

func (rep p2preply) Verify(origmd5 string) bool {
	return strings.Compare(rep.MD5, origmd5) == 0
}

func p2pPadding(origData []byte, blockSize ...int) []byte {
	var bk int
	if len(blockSize) > 0 {
		bk = blockSize[0]
	} else {
		bk = 128
	}
	padding := bk - len(origData)%bk
	padData := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(origData, padData...)
}

func p2pUnPadding(padData []byte) []byte {
	length := len(padData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(padData[length-1])
	return padData[:(length - unpadding)]
}

func ipfsapp() {
	// LibP2P code uses golog to log messages. They log with different
	// string IDs (i.e. "swarm"). We can control the verbosity level for
	// all loggers with:
	golog.SetAllLoggers(gologging.INFO) // Change to DEBUG for extra info

	// Parse options from the command line
	listenF := flag.Int("l", 8090, "wait for incoming connections")
	target := flag.String("d", "", "target peer to dial")
	secio := flag.Bool("secio", true, "enable secio")
	seed := flag.Int64("seed", 0, "set random seed for id generation")
	flag.Parse()

	if *listenF == 0 {
		log.Fatal("Please provide a port to bind on with -l")
	}

	// Make a host that listens on the given multiaddress
	has, err := makeBasicHosts(*listenF, *secio, *seed)
	if err != nil {
		log.Fatal(err)
	}
	for _, ha := range has {

		// Set a stream handler on host A. /echo/1.0.0 is
		// a user-defined protocol name.
		go func(ha host.Host) {
			ha.SetStreamHandler("/echo/1.0.0", func(s net.Stream) {
				log.Println("Got a new stream!")
				doEcho(s)

				s.Close()
				fmt.Println("Finish")
			})
		}(ha)
	}
	if *target == "" {
		log.Println("listening for connections")
		select {} // hang forever
	}
	/**** This is where the listener code ends ****/

	// The following code extracts target's the peer ID from the
	// given multiaddress
	ipfsaddr, err := ma.NewMultiaddr(*target)
	if err != nil {
		log.Fatalln(err)
	}

	pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
	if err != nil {
		log.Fatalln(err)
	}

	peerid, err := peer.IDB58Decode(pid)
	if err != nil {
		log.Fatalln(err)
	}

	// Decapsulate the /ipfs/<peerID> part from the target
	// /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
	targetPeerAddr, _ := ma.NewMultiaddr(
		fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerid)))
	targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)
	targetIP := strings.Split(targetAddr.String(), "/")[2]
	fmt.Println("I am target: ", targetIP)
	ip := "127.0.0.1"
	//for ip := range has {

	// We have a peer ID and a targetAddr so we add it to the peerstore
	// so LibP2P knows how to contact it
	has[ip].Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)

	log.Println("opening stream")
	// make a new stream from host B to host A
	// it should be handled on host A by the handler we set above because
	// we use the same /echo/1.0.0 protocol
	log.Printf(peerid.String())
	s, err := has[ip].NewStream(context.Background(), peerid, "/echo/1.0.0")
	defer s.Close()
	if err != nil {
		log.Fatalln(err)
	}
	args := flag.Args()
	r := ""
	if len(args) > 0 {
		r = args[0]
	}

	fmt.Println("r: ", r)
	arg := "cv.exe"
	t1 := time.Now()

	data, _ := ioutil.ReadFile(arg)

	req, _ := json.Marshal(&p2prequest{Op: "put", DataSize: len(data), HasPubKey: false})
	head := p2pPadding(req)
	data = append(head, data...)
	fmt.Println("Orign size: ", len(data))
	s.Write(data)
	fmt.Println("Finish Send")
	suffix := ""
	seq := strings.Split(arg, ".")
	if len(seq) > 1 {
		suffix += ("." + seq[len(seq)-1])
	}
	//ioutil.WriteFile("out"+suffix, data, 0644)
	//_, err = s.Write([]byte("Hello, world!\n"))

	if err != nil {
		log.Fatalln(err)
	}

	//out := []byte("HHH")
	//buf := bufio.NewReaderSize(s, 128)
	//fmt.Println(<-bb)
	/**
			out := make([]byte, 128)
			//n, err := s.Read(out)
			buf := bufio.NewReaderSize(s, 128)
			n, err := buf.Read(out)
			fmt.Println("n: ", n)
			fmt.Println("Reply size: ", len(out))
			if err != nil {
				log.Fatalln(err)
			}
	**/

	fmt.Println("???")
	head = make([]byte, 128)
	s.Read(head)
	//copy(head, out)
	head = PKCS5UnPadding(head)
	rep := &p2preply{}
	if err == nil {
		json.Unmarshal(head, rep)
		log.Printf("read reply: %s\n", string(head))
		log.Printf("Send Size: %d\n", rep.GetSize)
		log.Printf("Status: %v\n", rep.Status)
		log.Println("Verigy: ", rep.Verify(Md5SumBytes(data[128:])))
	}
	//ioutil.WriteFile("reply"+suffix, out, 0644)

	rt := time.Since(t1)
	fmt.Println("Transport Time: ", rt)
	//	}
}

// doEcho reads a line of data a stream and writes it back
func doEcho(s net.Stream) {
	var data []byte
	head := make([]byte, 128)
	n, err := s.Read(head)
	req := &p2prequest{}
	json.Unmarshal(PKCS5UnPadding(head), req)
	fmt.Println(string(PKCS5UnPadding(head)))
	if req.HasPubKey {
		pubkey := make([]byte, 1024)
		n, err = s.Read(pubkey)
		if err == nil && n == 1024 {
			req.PubKey = string(pubkey)
		} else {
			fmt.Println("Error at getting Public Key")
		}
	}
	b := make([]byte, mb/4)
	bs, _ := s.Read(b)
	data = append(data, b[:bs]...)
	tn := bs
	fmt.Println("bs: ", bs)
	for n, err := s.Read(b); err == nil; {
		//fmt.Printf("%d\n", n)
		data = append(data, b[:n]...)
		tn += n
		fmt.Println(n)
		fmt.Println(tn)
		if n < bs {
			break
		}
		n, err = s.Read(b)
		fmt.Println("After")
	}
	fmt.Println("Get Size: ", len(data))
	rep, _ := json.Marshal(&p2preply{Status: req.DataSize == len(data), Op: req.Op, GetSize: len(data), MD5: Md5SumBytes(data)})
	fmt.Println(len(p2pPadding(rep)))
	n, err = s.Write(p2pPadding(rep))
	fmt.Printf("%d : %v\n", n, err)
	fmt.Println("Reply")
	fmt.Println("Has PubKey? ", req.HasPubKey)
}

func appendToFile(fileName string, content string) error {
	// 以只写的模式，打开文件
	f, err := os.OpenFile(fileName, os.O_WRONLY, 0644)
	defer f.Close()
	if err != nil {
		fmt.Println("cacheFileList.yml file create failed. err: " + err.Error())
	} else {
		// 查找文件末尾的偏移量
		n, _ := f.Seek(0, os.SEEK_END)
		// 从末尾的偏移量开始写入内容
		_, err = f.WriteAt([]byte(content), n)
	}
	return err
}
func getMD5Sum(path string, isFile bool, merge ...bool) {
	//flag.Parse()
	t1 := time.Now()
	//返回Usage

	if isFile {
		result, err := Md5SumFile(path)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%x %s\n", result, path) //这里是*file
	} else {
		result, err := Md5SumFolder(path)
		if err != nil {
			panic(err)
		}
		m := false
		if len(merge) > 0 {
			m = merge[0]
		}
		// 开启merge，则只计算总的MD5
		if m {
			var s string
			for _, v := range result {
				s += fmt.Sprintf("%x", v)
			}
			fmt.Printf("%x %s\n", md5.Sum([]byte(s)), path)
		} else {
			for k, v := range result {
				fmt.Printf("%x %s\n", v, k)
			}
		}
	}
	rt := time.Since(t1)
	fmt.Println(rt)
}

func Md5SumFile(file string) (value [md5.Size]byte, err error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	value = md5.Sum(data)
	return value, nil
}

func Md5SumBytes(data []byte) string {
	value := md5.Sum(data)
	//md5val := make([]byte, 0, md5.Size)
	//md5val = append(md5val, value...)
	return fmt.Sprintf("%x", value)
}

func Md5SumFolder(folder string) (map[string][md5.Size]byte, error) {
	results := make(map[string][md5.Size]byte)
	filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		//判断文件属性
		if info.Mode().IsRegular() {
			result, err := Md5SumFile(path)
			if err != nil {
				return err
			}
			results[path] = result
		}
		return nil
	})
	return results, nil
}

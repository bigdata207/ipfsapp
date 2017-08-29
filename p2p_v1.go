package main

import (
	_ "bufio"
	"bytes"
	"context"
	"crypto/rand"
	"encoding/gob"
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
	gonet "net"
	"os"
	"strings"
	"time"
)

func getIP() []string {
	addrs, err := gonet.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ips := make([]string, 0)
	for _, address := range addrs {

		// ipnet.IP.IsLoopback()检查ip地址判断是否回环地址
		//if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
		if ipnet, ok := address.(*gonet.IPNet); ok {
			if ipnet.IP.To4() != nil {
				fmt.Println(ipnet.IP.String())
				ips = append(ips, ipnet.IP.String())
			}

		}
	}
	return ips
}

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
	ips := getIP()
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

var MB int = 1024 * 1024

type Request struct {
	Op       string
	DataSize int
}

func deepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}
func p2p(fname string) {
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
				defer s.Close()
				doEcho(s)
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

	t1 := time.Now()

	data, _ := ioutil.ReadFile(fname)
	fmt.Println("Orign size: ", len(data))
	s.Write(data)
	suffix := ""
	seq := strings.Split(fname, ".")
	if len(seq) > 1 {
		suffix += ("." + seq[len(seq)-1])
	}
	//ioutil.WriteFile("out"+suffix, data, 0644)
	//_, err = s.Write([]byte("Hello, world!\n"))

	if err != nil {
		log.Fatalln(err)
	}

	out, err := ioutil.ReadAll(s)
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
	log.Printf("read reply: %d\n", len(out))

	//ioutil.WriteFile("reply"+suffix, out, 0644)

	rt := time.Since(t1)
	fmt.Println("Transport Time: ", rt)
	//	}

}

// doEcho reads a line of data a stream and writes it back
func doEcho(s net.Stream) {
	var data []byte
	/**
		str, err := buf.ReadString('\n')
		for ; err == nil; str, err = buf.ReadString('\n') {
			if buf.Buffered() == 0 {
				if buf.Buffered() == 0 {
					break
				}
			}
			tol += str
			//fmt.Println("size: ", buf.Buffered())
		}
		if err == nil {
			tol += str
		}
	**/
	//tol += string(b[:n])
	//fmt.Println("b:", string(b[:n]))

	//}

	b := make([]byte, MB/4)
	bs, _ := s.Read(b)
	data = append(data, b[:bs]...)
	tn := bs
	//fmt.Println("bs: ", bs)
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
	/**
		bs := 262144
		b := make([]byte, bs)
		tn := 0
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
	**/
	fmt.Println("Get Size: ", len(data))
	s.Write(data)
	fmt.Println("Reply")
}

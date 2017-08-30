package ipfsapp

import (
	"flag"
	"fmt"
	_ "github.com/coreos/etcd/client"
	"github.com/rcrowley/go-metrics"
	"github.com/smallnest/rpcx"
	_ "github.com/smallnest/rpcx/log"
	"github.com/smallnest/rpcx/plugin"
	"testing"
	"time"
)

type Args struct {
	A int `msg:"a"`
	B int `msg:"b"`
}
type Reply struct {
	C int `msg:"c"`
}
type Arith int

var addr = flag.String("s", "127.0.0.1:8972", "service address")
var addr1 = flag.String("s1", "127.0.0.1:8973", "service address")

//var e = flag.String("e", "http://172.17.0.4:2379", "etcd URL")
var e = flag.String("e", "http://127.0.0.1:22379", "etcd URL")
var n = flag.String("n", "Arith", "Service Name")

func (a *Arith) Mul(args *Args, reply *Reply) error {
	reply.C = args.A * args.B
	return nil
}

type Arith1 int

func (t *Arith1) Mul(args *Args, reply *Reply) error {
	reply.C = args.A * args.B * 10
	return nil
}

func TestRpcxServer(t *testing.T) {
	fmt.Println("Here is Rpcx Server")
	flag.Parse()

	server := rpcx.NewServer()
	rplugin := &plugin.EtcdRegisterPlugin{
		ServiceAddress: "tcp@" + *addr,
		EtcdServers:    []string{*e},
		BasePath:       "/rpcx",
		Metrics:        metrics.NewRegistry(),
		Services:       make([]string, 1),
		UpdateInterval: time.Minute,
	}
	rplugin.Start()
	server.PluginContainer.Add(rplugin)
	server.PluginContainer.Add(plugin.NewMetricsPlugin())
	server.RegisterName(*n, new(Arith), "weight=1&m=devops")
	server.Serve("tcp", *addr)

	server1 := rpcx.NewServer()
	rplugin1 := &plugin.EtcdRegisterPlugin{
		ServiceAddress: "tcp@" + *addr1,
		EtcdServers:    []string{*e},
		BasePath:       "/rpcx",
		Metrics:        metrics.NewRegistry(),
		Services:       make([]string, 1),
		UpdateInterval: time.Minute,
	}
	rplugin1.Start()
	server1.PluginContainer.Add(rplugin1)
	server1.PluginContainer.Add(plugin.NewMetricsPlugin())
	server1.RegisterName(*n, new(Arith1), "weight=1.2&m=devops")
	server1.Serve("tcp", *addr1)
}

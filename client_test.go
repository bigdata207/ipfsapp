package ipfsapp

import (
	"context"
	"flag"
	"fmt"
	"github.com/smallnest/rpcx"
	"github.com/smallnest/rpcx/clientselector"
	"testing"
	"time"
)

//var e = flag.String("e", "http://127.0.0.1:2379", "etcd URL")

//var n = flag.String("n", "Arith", "Service Name")

func TestRpcxClient(t *testing.T) {
	fmt.Println("Here is Rpcx Client")
	/**
	fmt.Println("test")
	s := &rpcx.DirectClientSelector{Network: "tcp", Address: "127.0.0.1:8972", DialTimeout: 10 * time.Second}
	fmt.Println("test")
	client := rpcx.NewClient(s)
	defer client.Close()
	fmt.Println("test")
	args := &Args{7, 8}
	var reply Reply
	err := client.Call(context.Background(), "Arith.Mul", args, &reply)
	fmt.Println("test")
	if err != nil {
		fmt.Println("Error")
	}
	fmt.Printf("%d * %d = %d\n", args.A, args.B, reply.C)
	**/
	flag.Parse()

	//basePath = "/rpcx/" + serviceName

	for i := 0; i < 5; i++ {
		s := clientselector.NewEtcdClientSelector([]string{*e}, "/rpcx/"+*n, time.Minute, rpcx.RandomSelect, time.Minute)
		client := rpcx.NewClient(s)

		args := &Args{7, 8}
		var reply Reply
		err := client.Call(context.Background(), *n+".Mul", args, &reply)
		if err != nil {
			fmt.Printf("error for "+*n+": %d*%d, %v \n", args.A, args.B, err)
		} else {
			fmt.Printf(*n+": %d*%d=%d \n", args.A, args.B, reply.C)
		}
	}
}

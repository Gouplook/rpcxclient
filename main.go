/**
 * @Author: yinjinlin
 * @File:  main.go
 * @Description:
 * @Date: 2021/8/18 下午4:10
 */

package main

import (
	"context"
	"flag"
	example "github.com/rpcxio/rpcx-examples"
	"github.com/smallnest/rpcx/client"
	"log"
	"time"
)

var (
	addr = flag.String("addr", "localhost:8972", "server address")
)

// client
func main() {
	// Parse parses the command-line flags from os.Args[1:]. Must be called
	// after all flags are defined and before flags are accessed by the program
	flag.Parse()
	d, _ := client.NewPeer2PeerDiscovery("tcp"+*addr, "")
	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()

	args := &example.Args{
		A: 10,
		B: 20,

	}

	for {
		reply := &example.Reply{}
		err := xclient.Call(context.Background(),"Mul", args, reply)
		if err != nil {
			log.Fatalf("failed to call: %v",err)
		}
		log.Printf("%d * %d = %d", args.A, args.B, reply.C)
		replyAdd := &example.Reply{}
		err = xclient.Call(context.Background(),"Add", args, replyAdd)
		if err != nil {
			log.Fatalf("%d * %d = %d",args.A, args.B,replyAdd.C)
		}

		time.Sleep(1e9)
	}

}

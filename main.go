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
	"github.com/smallnest/rpcx/server"
	"log"
	"time"
)

var (
	addr     = flag.String("addr", "localhost:10002", "server address")
	zkAddr   = flag.String("zkAddr", "localhost:2181", "zookeeper address")
	basePath = flag.String("base", "/rpc_test", "prefix path")
)

// client
func main() {
	flag.Parse()
	// d, _ := client.NewPeer2PeerDiscovery("tcp@"+*addr, "")
	d, _ := client.NewZookeeperDiscovery(*basePath, "Arith", []string{*zkAddr}, nil)
	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()
	// IsClosing 表示客户端是关闭着的，并且不会接受新的调用。

	args := &example.Args{
		A: 10,
		B: 20,
	}

	for {
		reply := &example.Reply{}
		// 调用乘法
		err := xclient.Call(context.Background(), "Mul", args, reply)
		if err != nil {
			// 相当于printf（但是有exit(1) 退出）
			log.Fatalf("failed to call: %v", err)
		}
		log.Printf("%d * %d = %d", args.A, args.B, reply.C)

		// 调用加法
		replyAdd := &example.Reply{}
		err = xclient.Call(context.Background(), "Add", args, replyAdd)
		if err != nil {
			log.Fatalf("failed to call: %v", err)
		}
		log.Printf("%d + %d = %d", args.A, args.B, replyAdd.C)

		time.Sleep(1e9)
	}

}

// 异步请求服务
var (
	// asyAdrr = flag.String("addr", "localhost:10002", "service address")
	asyAdrr = "xx"
)

func asynchronism() {
	flag.Parse()
	// 元数据：不是服务请求和服务响应数据的业务数据，而是一些辅助性的数据。
	d, _ := client.NewPeer2PeerDiscovery("tcp@"+asyAdrr, "")

	//
	// d := client.NewInprocessDiscovery()
	// Failtry：选择相同节点并重试，直到达到最大重试次数
	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)

	func() {
		defer xclient.Close()
	}()


	args := &example.Args{
		A: 100,
		B: 200,
	}

	for i := 0; i < 100; i++ {
		reply := &example.Reply{}
		call, err := xclient.Go(context.Background(), "Mul", args, reply, nil)
		if err != nil {
			log.Fatalf("failed to call: %v", err)
		}
		replyCall := <-call.Done
		if replyCall.Error != nil {
			log.Fatalf("failed to call:%v", replyCall.Error)
		} else {
			log.Fatalf("%d * %d", args.A, args.B, reply.C)
		}
	}



}

func addRegistryPlugin(s *server.Server){
	r := client.InprocessClient
	s.Plugins.Add(r)
}

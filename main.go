package main

import (
	"net/rpc"
	"net"
	"log"
	"fmt"
	"net/http"
	"SimpleRaft/settings"
	"SimpleRaft/servers"
)

func StartService() {
	m := new(servers.RaftManager)
	m.Init()

	rrpc := &servers.RaftRPC{m}
	rpc.Register(rrpc)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", m.CurrentIP+":"+settings.SERVERPORT)
	if e != nil {
		log.Fatal(e)
	}
	fmt.Println("启动RPC监听服务...")
	http.Serve(l, nil)
}
func main() {
	StartService()
}

package main

import (
	"net/rpc"
	"net"
	"log"
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
	log.Println("启动RPC监听服务...")
	http.Serve(l, nil)
}

func f(dct map[string]string)  {
	dct["hello"] = "CES"
}

func main() {

}

package main

import (
	"net/rpc"
	"net"
	"log"
	"fmt"
	"net/http"
	"SimpleRaft/settings"
)

func StartService()  {
	m := new(RaftManager)
	m.StartRoleService()
	rpc.Register(m)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", m.CurrentIP + ":" + settings.SERVERPORT)
	if e != nil {
		log.Fatal(e)
	}
	fmt.Println("启动服务...")
	http.Serve(l, nil)
}
func main() {
	StartService()
}

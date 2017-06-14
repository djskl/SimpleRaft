package main

import (
	"fmt"
	"net/rpc"
	"net"
	"log"
	"net/http"
)

func StartService()  {
	m := new(RaftManager)
	m.StartRoleService()
	rpc.Register(m)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", "127.0.0.1:1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	fmt.Println("I'm listening...")
	http.Serve(l, nil)
}

func main() {
	StartService()
}

package main

import (
	"fmt"
	"time"
)

//func StartService()  {
//	m := RaftManager{}
//	rpc.Register(m)
//	rpc.HandleHTTP()
//	l, e := net.Listen("tcp", "127.0.0.1:1234")
//	if e != nil {
//		log.Fatal("listen error:", e)
//	}
//	fmt.Println("I'm listening...")
//	http.Serve(l, nil)
//}

func main() {
	c := time.Tick(1 * time.Second)
	for now := range c {
		fmt.Printf("%v\n", now)
	}
}

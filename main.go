package main

import (
	"fmt"
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

	var ch chan int

	if ch == nil{
		fmt.Println("hello")
	}

	//go func() {
	//	ch <- 3
	//}()
	//
	//x := <-ch
	//
	//fmt.Println(x)
}

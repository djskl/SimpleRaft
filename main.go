package main

import (
	"fmt"
	"math/rand"
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

func f(x []int, y []int) []int {
	x = append(x, y...)
	return x
}

func main() {
	for idx :=10;idx<100;idx++ {
		fmt.Println(rand.Intn(3000) + 3000)
	}

}

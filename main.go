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

func f(x []int, y []int) []int {
	x = append(x, y...)
	return x
}

func main() {

	x := 10
	y := 3

	if ( x > 5 ) || ( y > 5 ) {
		fmt.Println("hello")
	}

}

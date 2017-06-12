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

func f(x []int, y []int) []int {
	x = append(x, y...)
	return x
}

func main() {
	ch_ot := make(chan bool)
	go func() {
		for {
			select {
			case x := <-ch_ot:
				fmt.Println(x)
			case <-time.After(200 * time.Millisecond):
				fmt.Println("timeout")
			}
		}
	}()

	go func() {
		for {
			ch_ot <- true
			time.Sleep(time.Millisecond*100)
		}
	}()

	time.Sleep(time.Second*10)

}

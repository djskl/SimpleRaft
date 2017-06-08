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

func main()  {
	dct := map[string]string{
		"one": "1",
		"two": "2",
	}

	dct["six"] = "6"

	for k, v := range dct {
		fmt.Println(k, v)
	}

}
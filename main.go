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

type TypeError struct {
	s string
}
func (t TypeError) Error() string {
	return t.s
}

func getBirthday(age interface{}) (int, error) {
	val, ok := age.(int)
	if(!ok){
		return 0, TypeError{"年龄必须是整数"}
	}
	return 2017 - val, nil
}
func main()  {
	age := "hello"
	_, err := getBirthday(age)
	if(err != nil){
		val, ok := err.(TypeError)
		fmt.Println(val, ok)
	}
}
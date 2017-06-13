package main

import (
	"sync"
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

type CountWG struct {
	wg sync.WaitGroup
	num int
	wk sync.Mutex
}

func (this *CountWG) Add(delta int) {
	this.wk.Lock()
	defer this.wk.Unlock()
	this.num += delta
	this.wg.Add(delta)
}

func (this *CountWG) Done() {
	this.wk.Lock()
	defer this.wk.Unlock()
	if this.num == 0 {
		return
	}
	this.wg.Done()
	this.num -= 1
}

func (this *CountWG) Size() int {
	return this.num
}

func (this *CountWG) Wait() {
	this.wg.Wait()
}

func main() {
	wg := new(CountWG)
	wg.Add(1)

	go func() {
		defer wg.Done()
		wg.Add(1)
		fmt.Println("1")
		time.Sleep(time.Millisecond*100)
		go func() {
			defer wg.Done()
			wg.Add(2)
			fmt.Println("2")
			time.Sleep(time.Millisecond*100)

			go func() {
				defer wg.Done()
				fmt.Println("hello")
				time.Sleep(time.Millisecond*500)
				fmt.Println("world")
			}()

			go func() {
				defer wg.Done()
				wg.Add(1)
				fmt.Println("3")
				time.Sleep(time.Millisecond*100)
				go func() {
					defer wg.Done()
					wg.Add(1)
					fmt.Println("4")
					time.Sleep(time.Millisecond*100)
					go func() {
						defer wg.Done()
						wg.Add(1)
						fmt.Println("5")
						time.Sleep(time.Millisecond*100)
						go func() {
							defer wg.Done()
							wg.Add(1)
							fmt.Println("6")
							time.Sleep(time.Millisecond*100)
							go func() {
								defer wg.Done()
								wg.Add(1)
								fmt.Println("7")
								time.Sleep(time.Millisecond*100)
								go func() {
									defer wg.Done()
									wg.Add(1)
									fmt.Println("8")
									time.Sleep(time.Millisecond*100)
									go func() {
										defer wg.Done()
										wg.Add(1)
										fmt.Println("9")
										time.Sleep(time.Millisecond*100)
									}()
								}()
							}()
						}()
					}()
				}()
			}()
		}()
	}()

	go func() {
		time.Sleep(time.Millisecond*500)
		var x int
		for {
			if x=wg.Size();x>0{
				wg.Done()
			}else{
				break
			}
		}
	}()

	wg.Wait()

}

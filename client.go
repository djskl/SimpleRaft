package main

import (
	"net/rpc"
	"log"
	"sync"
	"math/rand"
	"strconv"
	"time"
)

var ALLSERVERS = [5]string{"192.168.16.2", "192.168.16.3", "192.168.16.4", "192.168.16.5", "192.168.16.6"}

var LEADER_IP string = ""

type CommandAck struct {
	Ok       bool
	Cmd      string
	LeaderIP string
}

func Submit(cmd string) {
	var client *rpc.Client
	var err error
	if LEADER_IP == "" {
		for idx := 0; idx < 5; idx++ {
			server_ip := ALLSERVERS[idx]
			client, err = rpc.DialHTTP("tcp", server_ip+":5656")
			if err != nil {
				log.Printf("无法与%s建立连接!!!\n", server_ip)
			} else {
				break
			}
		}
	} else {
		client, err = rpc.DialHTTP("tcp", LEADER_IP+":5656")
		if err != nil {
			log.Printf("无法与%s建立连接!!!\n", LEADER_IP)
			LEADER_IP = ""
			Submit(cmd)
			return
		}
	}

	if client == nil {
		log.Println("无服务器可用!!!")
		return
	}

	ack := new(CommandAck)
	err = client.Call("RaftRPC.Command", cmd, ack)
	if err != nil {
		log.Println("命令提交失败：", err)
		return
	}

	client.Close()

	if ack.Ok {
		log.Printf("命令：%s提交成功!!!\n", cmd)
	} else {
		if ack.LeaderIP != "" {
			log.Printf("LEADER_IP是：%s，重新提交...\n", ack.LeaderIP)
			LEADER_IP = ack.LeaderIP
			Submit(cmd)
		} else {
			log.Printf("命令：%s提交失败!!!\n", cmd)
		}
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(100)
	for idx := 0; idx < 100; idx++ {
		go func() {
			defer wg.Done()
			for idy:=0;idy<100;idy++{
				n := strconv.Itoa(rand.Intn(10))
				Submit(n)
				time.Sleep(time.Duration(500 + rand.Intn(500))*time.Millisecond)
			}
		}()
	}
	wg.Wait()
}

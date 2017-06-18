package main

import (
	"net/rpc"
	"log"
)

var ALLSERVERS = [5]string{"10.0.138.151", "10.0.138.152", "10.0.138.153", "10.0.138.155", "10.0.138.158"}

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
		if ack.Cmd == cmd {
			log.Printf("命令：%s提交成功!!!\n", cmd)
		} else {
			log.Printf("命令：%s被覆盖(%s)!!!\n", cmd, ack.Cmd)
		}
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
	ALLCMDS := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	for _, v := range ALLCMDS {
		Submit(v)
	}
}

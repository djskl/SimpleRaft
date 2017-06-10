package servers

import (
	"SimpleRaft/settings"
	"log"
	"net/rpc"
	"SimpleRaft/utils"
)

type Leader struct {
	*BaseRole
	NextIndex  map[string]int
	MatchIndex map[string]int

	chan_newlog chan int //当有新日志到来时激活日志复制服务
	log_commits []int    //index位置的日志已复制成功的机器数量

	chan_clients []chan int	//leader确定提交了某项日志后，激活这一组管道
	chan_newlogs chan int	//客户端添加新日志后，激活这个管道
}

func (this *Leader) Init() error {
	this.BaseRole.init()	//相当于调用父类的构造函数
	this.chan_newlog = make(chan int)
	this.MatchIndex = make(map[string]int)
	this.NextIndex = make(map[string]int)

	log_length := len(this.Logs)
	for _, ip := range settings.AllServers {
		this.NextIndex[ip] = log_length
	}
	return nil
}

//启动日志replicate服务
func (this *Leader) startReplService() {
	for {
		log_length := len(this.Logs)
		for idx := 0; idx < len(settings.AllServers); idx++ {
			ip := settings.AllServers[idx]
			if ip == this.IP {
				continue
			}
			next_index := this.NextIndex[ip]
			if next_index > log_length {
				continue
			}
			go this.replicateLog(ip, next_index)
		}
	}
}

func (this *Leader) replicateLog(ip string, next_index int) {
	preLog := this.Logs[next_index-1]
	toSendEntries := this.Logs[next_index:]

	logReq := LogAppArg{
		Term:            this.CurrentTerm,
		LeaderID:        this.IP,
		LeaderCommitIdx: this.CommitIndex,
		PreLogIndex:     next_index - 1,
		PreLogTerm:      preLog.Term,
		Entries:         toSendEntries,
	}

	logAck := new(LogAckArg)
	for {	//如果连接(或RPC调用)失败，就反复尝试
		client, err := rpc.DialHTTP("tcp", ip+":"+settings.SERVERPORT)
		if err != nil {
			log.Printf("failed to connect %s from %s\n", ip, this.IP)
			continue
		}

		err = client.Call("RaftManager.AppendLog", logReq, logAck)
		if err != nil {
			log.Printf("appendLog failed: %s--->%s", this.IP, ip)
			continue
		}
		break
	}
	this.handleLogAck(logAck)
}

// 1、统计提交成功的次数(commits>=3提交，更新commitIdx时，激活所有正在等待的客户端)；
// 2、切换角色(currentTerm<followerTerm);
func (this *Leader) handleLogAck(logAck *LogAckArg) {

	//当前这个leader已经过期了，通知manager将当前角色改为Follower
	if logAck.Term > this.CurrentTerm {
		this.CurrentTerm = logAck.Term
		this.SetAlive(false)
		this.chan_role <- settings.FOLLOWER	//告诉manager：将当前角色切换为Follower
		return
	}



}

func (this *Leader) HandleVoteReq(args0 VoteReqArg, args1 *VoteAckArg) error {
	return nil
}

func (this *Leader) HandleAppendLogReq(args0 LogAppArg, args1 *LogAckArg) error {
	return nil
}
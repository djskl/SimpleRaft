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

	chan_newlog chan int 	//当有新日志到来时激活日志复制服务
	log_commits []int		//对应日志已复制成功的机器数量
}

func (this *Leader) Init() error {
	this.IsAlive = true
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
	nextLogIdx := this.NextIndex[ip]
	preLog := this.Logs[nextLogIdx-1]
	toSendEntries := this.Logs[nextLogIdx:]

	inArg := LogAppArg{
		Term:            this.CurrentTerm,
		LeaderID:        this.IP,
		LeaderCommitIdx: this.CommitIndex,
		PreLogIndex:     nextLogIdx - 1,
		PreLogTerm:      preLog.Term,
		Entries:         toSendEntries,
	}
	client, err := rpc.DialHTTP("tcp", ip+":"+settings.SERVERPORT)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	err = client.Call("RaftManager.AppendLog", inArg, ack)
	if err != nil {
		log.Fatalf("RaftManager AppendLog Error: NextLogIndex: %d (LogTerm: %d)",
			nextLogIdx, this.CurrentTerm)
	}

}

func (this *Leader) HandleVoteReq(args0 VoteReqArg, args1 *VoteAckArg) error {
	return nil
}

func (this *Leader) HandleAppendLogReq(args0 LogAppArg, args1 *LogAckArg) error {
	return nil
}

func (this *Leader) SetAlive(alive bool) {
	this.IsAlive = alive
}

//将用户发出的命令发送到各个机器
func (this *Leader) HandleCommandReq(cmds string, success *bool) error {

	var cmts int

	ok_chan := make(chan bool)

	for idx := 0; idx < len(settings.AllServers); idx++ {
		ip := settings.AllServers[idx]
		if ip == this.IP {
			continue
		}
		go func(ip string) {
			var rst *LogAckArg
			rst = this.replicateLog(ip)
			for {
				if rst.Term > this.CurrentTerm { //抛出异常，由RaftManager将当前角色设置为Follower
					ok_chan <- false
					break
				}

				if !rst.Success { //因为数据不一致导致数据添加失败(更新nextIndex,然后重试)
					this.NextIndex[ip] -= 1
					rst = this.replicateLog(ip)
				} else {
					ok_chan <- true
					break
				}
			}
		}(ip)

		for {
			isOK := <-ok_chan
			if isOK {
				if cmts++; cmts > 1 {
					*success = true
					return nil
				}
			} else {
				*success = false
				return utils.TermError{"the leader is expired"}
			}
		}
	}
	return nil
}

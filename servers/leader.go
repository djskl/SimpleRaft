package servers

import (
	"SimpleRaft/settings"
	"log"
	"net/rpc"
	"SimpleRaft/db"
	"SimpleRaft/clog"
)

type Leader struct {
	*BaseRole
	NextIndex  map[string]int
	MatchIndex map[string]int

	chan_newlog chan int //当有新日志到来时激活日志复制服务
	log_commits []int    //index位置的日志已复制成功的机器数量

	chan_clients []chan int //leader确定提交了某项日志后，激活这一组管道
	chan_newlogs chan bool  //客户端添加新日志后，激活这个管道
	chan_commits chan int   //更新了commit后，激活这个管道
}

func (this *Leader) Init(role_chan chan int) error {
	this.BaseRole.init(role_chan) //相当于调用父类的构造函数
	this.chan_newlog = make(chan int)
	this.MatchIndex = make(map[string]int)
	this.NextIndex = make(map[string]int)

	log_length := this.Logs.Size()
	for _, ip := range settings.AllServers {
		this.NextIndex[ip] = log_length
	}
	return nil
}

func (this *Leader) startAllService() {
	go this.startLogReplService()
	go this.startLogApplService()
}

//日志replicate服务
func (this *Leader) startLogReplService() {
	for {
		log_length := this.Logs.Size()
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

		<-this.chan_newlogs

	}
}

//日志应用服务(更新lastApplied)
func (this *Leader) startLogApplService() {
	for {
		for this.LastApplied < this.CommitIndex {
			_log := this.Logs.Get(this.LastApplied+1)
			db.WriteToDisk(_log.Command)
			this.LastApplied++
		}
		<-this.chan_commits
	}
}

func (this *Leader) replicateLog(ip string, next_index int) {
	if next_index < 0 {
		return
	}

	var preidx, pretem int

	if next_index > 0 {
		preLog := this.Logs.Get(next_index-1)
		preidx = next_index - 1
		pretem = preLog.Term
	}

	toSendEntries := this.Logs.GetFrom(next_index)

	logReq := LogAppArg{
		Term:            this.CurrentTerm,
		LeaderID:        this.IP,
		LeaderCommitIdx: this.CommitIndex,
		PreLogIndex:     preidx,
		PreLogTerm:      pretem,
		Entries:         toSendEntries,
	}

	logAck := new(LogAckArg)
	for { //如果连接(或RPC调用)失败，就反复尝试
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
	this.handleLogAck(ip, next_index, logAck)
}

func (this *Leader) handleLogAck(ip string, next_index int, logAck *LogAckArg) {
	//处理leader过期的情况
	if logAck.Term > this.CurrentTerm {
		this.CurrentTerm = logAck.Term
		this.SetAlive(false)
		this.chan_role <- settings.FOLLOWER //告诉manager：将当前角色切换为Follower
		return
	}

	//数据不一致导致更新失败
	if !logAck.Success {
		go this.replicateLog(ip, next_index-1)
		return
	}

	this.NextIndex[ip] = logAck.LastLogIndex + 1
	this.MatchIndex[ip] = logAck.LastLogIndex

	this.updateCommitIndex(logAck.LastLogIndex)

}

// If there exists an N such that N > commitIndex,
// a majority of matchIndex[i] ≥ N and log[N].term == currentTerm:
// set commitIndex = N
func (this *Leader) updateCommitIndex(logIndex int) {
	if logIndex <= this.CommitIndex {
		return
	}

	for logIndex > this.CommitIndex {
		_log := this.Logs.Get(logIndex)

		if _log.Term < this.CurrentTerm { //不能提交上一个term的日志
			return
		}

		nums := 0
		for _, match_idx := range this.MatchIndex {
			if match_idx >= logIndex {
				if nums++; nums > 2 {
					this.CommitIndex = match_idx
					this.chan_commits <- this.CommitIndex
					return
				}
			}
		}
		logIndex--
	}
}

func (this *Leader) HandleCommandReq(cmd string, ok *bool) error {
	_log := clog.Item{this.CurrentTerm, cmd}
	this.Logs.Add(_log)

	go func() {
		this.chan_newlogs <- true
	}()

	return nil
}

func (this *Leader) HandleVoteReq(args0 VoteReqArg, args1 *VoteAckArg) error {
	return nil
}

func (this *Leader) HandleAppendLogReq(args0 LogAppArg, args1 *LogAckArg) error {
	return nil
}

package servers

import (
	"SimpleRaft/settings"
	"log"
	"net/rpc"
	"SimpleRaft/db"
	"SimpleRaft/clog"
	"time"
	"SimpleRaft/utils"
)

type Leader struct {
	*BaseRole

	//persistent state
	CurrentTerm int
	VotedFor    string

	AllServers []string
	
	NextIndex  map[string]int
	MatchIndex map[string]int

	Clients *ClientWaitGroup

	chan_newlog  chan int         //当有新日志到来时激活日志复制服务
	chan_commits chan int         //更新了commit后，激活这个管道

	//角色是否处于激活状态，供角色启动的子协程参考
	//用指针防止copy
	active *utils.AtomicBool
	chan_role *chan RoleState
}

func (this *Leader) Init() error {
	tmp := make(chan RoleState)
	this.chan_role = &tmp

	this.VotedFor = this.IP

	this.active = new(utils.AtomicBool)
	this.active.Set()

	this.chan_newlog = make(chan int)
	this.chan_commits = make(chan int)

	this.MatchIndex = make(map[string]int)
	this.NextIndex = make(map[string]int)

	log_length := this.Logs.Size()
	for _, ip := range this.AllServers {
		this.NextIndex[ip] = log_length + 1
	}

	log.Printf("LEADER(%d)：初始化...\n", this.CurrentTerm)

	return nil
}

func (this *Leader) SetAlive(alive bool){
	this.active.SetTo(alive)
}

func (this *Leader) GetAlive() bool {
	return this.active.IsSet()
}

func (this *Leader) GetRoleChan() *chan RoleState {
	return this.chan_role
}

//启动所有服务
func (this *Leader) StartAllService() {
	go this.startLogReplService()
	go this.startLogApplService()
}

func (this *Leader) HandleCommandReq(cmd string, cmdAck *CommandAck) error {
	if !this.active.IsSet() {
		return nil
	}

	log.Printf("LEADER(%d)：收到客户端请求：%s\n", this.CurrentTerm, cmd)

	_log := clog.LogItem{this.CurrentTerm, cmd}

	pos := this.Logs.Add(_log)	//pos从1开始计数，第一个日志的pos为1

	this.chan_newlog <- pos

	chan_client := make(chan int)
	this.Clients.Add(pos, chan_client)
	for {
		if !this.active.IsSet() {
			cmdAck.Ok = false
			cmdAck.LeaderIP = this.VotedFor
			break
		}

		select {
		case commitIdx := <-chan_client:
			if commitIdx >= pos {
				cmdAck.Ok = true
				cmdAck.LeaderIP = this.VotedFor
				cmdAck.Cmd = this.Logs.Get(pos).Command
				this.Clients.Remove(pos)
				log.Printf("LEADER(%d)：成功处理：%s\n", this.CurrentTerm, cmd)
				return nil
			}
		case <-time.After(time.Millisecond * time.Duration(settings.CLIENT_WAIT)):
			break	//跳出select循环，但不会跳出for循环
		}

	}
	log.Printf("LEADER(%d)：处理失败：%s\n", this.CurrentTerm, cmd)
	return nil
}

func (this *Leader) HandleVoteReq(args0 VoteReqArg, args1 *VoteAckArg) error {

	log.Printf("LEADER(%d)：收到来在%s的投票请求...\n", this.CurrentTerm, args0.CandidateID)

	args1.Term = this.CurrentTerm
	args1.VoteGranted = false

	//收到了比自己大的term直接转换为follower
	if args0.Term > this.CurrentTerm && this.active.IsSet() {
		log.Printf("LEADER(%d)：过期了，转为Follower...\n", this.CurrentTerm)
		this.CurrentTerm = args0.Term
		rolestate := RoleState{settings.FOLLOWER, this.CurrentTerm}
		select {
		case *this.chan_role <- rolestate:
			return nil
		case <-time.After(time.Millisecond * time.Duration(settings.COMMIT_WAIT)):
			return nil
		}
	}

	return nil
}

func (this *Leader) HandleAppendLogReq(args0 LogAppArg, args1 *LogAckArg) error {
	log.Printf("LEADER(%d)：收到append_log(TERM: %d)请求...\n", this.CurrentTerm, args0.Term)
	args1.Term = this.CurrentTerm
	args1.Success = false

	//收到了比自己大的term直接转换为follower
	if args0.Term > this.CurrentTerm && this.active.IsSet() {
		log.Printf("LEADER(%d)：过期了，转为Follower...\n", this.CurrentTerm)
		this.CurrentTerm = args0.Term
		rolestate := RoleState{settings.FOLLOWER, this.CurrentTerm}
		select {
		case *this.chan_role <- rolestate:
			return nil
		case <-time.After(time.Millisecond * time.Duration(settings.COMMIT_WAIT)):
			return nil
		}
	}
	return nil
}

//日志replicate服务：有日志传日志，没日志传心跳
func (this *Leader) startLogReplService() {
	log.Printf("LEADER(%d)：启动replicate_log服务...\n", this.CurrentTerm)
	for {
		if !this.active.IsSet() {
			break
		}

		for idx := 0; idx < len(this.AllServers); idx++ {
			ip := this.AllServers[idx]
			if ip == this.IP {
				continue
			}
			next_index := this.NextIndex[ip]
			go this.replicateLog(ip, next_index)
		}

		select {
		case <-this.chan_newlog:
			//do nothing
		case <-time.After(time.Millisecond*time.Duration(settings.HEART_BEATS)):
			//do nothing
		}
	}
	log.Printf("LEADER(%d)：replicate_log服务终止！！！\n", this.CurrentTerm)
}

//日志应用服务(更新lastApplied)
func (this *Leader) startLogApplService() {
	log.Printf("LEADER(%d)：启动日志应用服务...\n", this.CurrentTerm)
	for {
		if !this.active.IsSet() {
			break
		}

		for this.LastApplied <= this.CommitIndex {
			this.LastApplied++
			_log := this.Logs.Get(this.LastApplied)
			if _log.Command != "" {
				db.WriteToDisk(_log.Command)
				log.Printf("FOLLOWER(%d): %d 已写到日志文件\n", this.CurrentTerm, this.LastApplied)
			}
		}
		select {
		case <-this.chan_commits:
			//do nothing
		case <-time.After(time.Millisecond*time.Duration(settings.COMMIT_WAIT)):
			//do nothing
		}
	}
	log.Printf("LEADER(%d)：日志应用服务终止！！！\n", this.CurrentTerm)
}

func (this *Leader) replicateLog(ip string, next_index int) {
	if !this.active.IsSet() {
		return
	}

	if ip == "" || next_index < 1 {
		return
	}

	var preidx, pretem int

	if next_index == 1 { //此时复制第一个日志
		preidx = 0
	} else {
		preLog := this.Logs.Get(next_index - 1)
		preidx = next_index - 1
		pretem = preLog.Term
	}

	toSendEntries := this.Logs.GetFrom(next_index)
	if toSendEntries != nil && len(toSendEntries) > 0 {
		log.Printf("LEADER(%d)：向%s复制日志:%d...\n", this.CurrentTerm, ip, next_index)
	} else {
		//log.Printf("LEADER(%d)：向%s发送心跳信息...\n", this.CurrentTerm, ip)
	}

	logReq := LogAppArg{
		Term:            this.CurrentTerm,
		LeaderID:        this.IP,
		LeaderCommitIdx: this.CommitIndex,
		PreLogIndex:     preidx,
		PreLogTerm:      pretem,
		Entries:         toSendEntries,
	}

	logAck := new(LogAckArg)

	var client *rpc.Client
	var err error
	for {
		if !this.active.IsSet() {
			return
		}
		client, err = rpc.DialHTTP("tcp", ip+":"+settings.SERVERPORT)
		if err != nil {
			log.Printf("LEADER(%d)：无法与Follower(%s)建立数据连接！！！\n", this.CurrentTerm, ip)
			time.Sleep(time.Millisecond*time.Duration(settings.RPC_WAIT))
			continue
		}
		break
	}

	for {
		if !this.active.IsSet() {
			client.Close()
			return
		}
		err = client.Call("RaftRPC.AppendLog", logReq, logAck)
		if err != nil {
			log.Printf("LEADER(%d)：调用Follower(%s)的AppendLog方法失败！！！\n", this.CurrentTerm, ip)
			time.Sleep(time.Millisecond*time.Duration(settings.RPC_WAIT))
			continue
		}
		break
	}
	client.Close()
	this.handleLogAck(ip, next_index, logAck)
}

func (this *Leader) handleLogAck(ip string, next_index int, logAck *LogAckArg) {
	if !this.active.IsSet() {
		return
	}

	//处理leader过期的情况
	if logAck.Term > this.CurrentTerm && this.active.IsSet() {
		log.Printf("LEADER(%d)：过期了，转为Follower...\n", this.CurrentTerm)
		this.CurrentTerm = logAck.Term
		rolestate := RoleState{settings.FOLLOWER, this.CurrentTerm}
		select {
		case *this.chan_role <- rolestate:
			return
		case <-time.After(time.Millisecond * time.Duration(settings.CHANN_WAIT)):
			return
		}
	}

	//数据不一致，更新next_index重新replicate
	if !logAck.Success {
		log.Printf("LEADER(%d)：与Follower(%s)产生冲突，更新index重新replicate...\n", this.CurrentTerm, ip)
		go this.replicateLog(ip, next_index - 1)
		return
	}

	if next_index == logAck.LastLogIndex + 1 {
		//log.Printf("LEADER(%d)：收到Follower(%s)的心跳响应！！！\n", this.CurrentTerm, ip)
		return
	}

	this.NextIndex[ip] = logAck.LastLogIndex + 1
	this.MatchIndex[ip] = logAck.LastLogIndex

	this.updateCommitIndex(logAck.LastLogIndex)

}

// If there exists an N such that N > commitIndex,
// a majority of matchIndex[i] ≥ N and log[N].term == currentTerm:
// set commitIndex = N
func (this *Leader) updateCommitIndex(newLogIndex int) {
	if !this.active.IsSet() {
		return
	}

	//if newLogIndex <= this.CommitIndex {
	//	log.Printf("LEADER(%d)：日志(%d)已提交!!!\n", this.CurrentTerm, newLogIndex)
	//	go func() {
	//		this.chan_commits <- this.CommitIndex
	//		this.Clients.NotifyAll(this.CommitIndex)
	//	}()
	//	return
	//}

	for newLogIndex > this.CommitIndex {
		_log := this.Logs.Get(newLogIndex)

		if _log.Term < this.CurrentTerm { //不能提交上一个term的日志
			continue
		}

		nums := 1
		for _, match_idx := range this.MatchIndex {
			if match_idx >= newLogIndex {
				if nums++; nums >= settings.MAJORITY {
					this.CommitIndex = match_idx
					log.Printf("LEADER(%d)：日志复制成功，当前commitIndex=%d\n", this.CurrentTerm, this.CommitIndex)
					go func() {
						this.chan_commits <- this.CommitIndex
						this.Clients.NotifyAll(this.CommitIndex)
					}()
					return
				}
			}
		}
		newLogIndex--
	}
}

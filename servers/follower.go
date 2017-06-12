package servers

import "SimpleRaft/db"

type Follower struct {
	*BaseRole
	chan_commits chan int //更新了commit后，激活这个管道
}

func (this *Follower) Init(role_chan chan int) error {
	this.BaseRole.init(role_chan) //相当于调用父类的构造函数
	this.chan_commits = make(chan int)
	return nil
}

func (this *Follower) StartAllService() {

}

func (this *Follower) HandleVoteReq(args0 VoteReqArg, args1 *VoteAckArg) error {
	return nil
}

func (this *Follower) HandleAppendLogReq(args0 LogAppArg, args1 *LogAckArg) error {
	//收到了过期leader的log_rpc
	if args0.Term > this.CurrentTerm {
		args1.Term = this.CurrentTerm
		args1.Success = false
		return nil
	}

	this.CurrentTerm = args1.Term

	//来自leader的心跳信息(更新commit index)
	if len(args0.Entries) == 0 {
		return this.handleHeartBeat(args0, args1)
	}

	args1.Term = this.CurrentTerm
	var lastLogIdx int = -1
	if args0.PreLogIndex > -1 {
		preLog := this.Logs.Get(args0.PreLogIndex)
		if preLog.Term != args0.PreLogTerm {
			this.Logs.RemoveFrom(args0.PreLogIndex)
			args1.Success = false
			return nil
		}
	}
	lastLogIdx = this.Logs.Extend(args0.Entries)
	args1.LastLogIndex = lastLogIdx
	args1.Success = true

	return nil
}

func (this *Follower) HandleCommandReq(cmd string, ok *bool) error {
	return nil
}

func (this *Follower) handleHeartBeat(args0 LogAppArg, args1 *LogAckArg) error {
	if args0.LeaderCommitIdx > this.CommitIndex {
		log_size := this.Logs.Size()
		if args0.LeaderCommitIdx < log_size {
			this.CommitIndex = args0.LeaderCommitIdx
		} else {
			this.CommitIndex = log_size
		}
		go func() {
			this.chan_commits <- this.CommitIndex
		}()
	}
	return nil
}

//日志应用服务(更新lastApplied)
func (this *Follower) startLogApplService() {
	for {
		for this.LastApplied < this.CommitIndex {
			_log := this.Logs.Get(this.LastApplied + 1)
			db.WriteToDisk(_log.Command)
			this.LastApplied++
		}
		<-this.chan_commits
	}
}

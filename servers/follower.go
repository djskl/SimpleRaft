package servers

import (
	"SimpleRaft/db"
	"time"
	"SimpleRaft/settings"
	"math/rand"
	"SimpleRaft/utils"
)

type Follower struct {
	*BaseRole
	//persistent state
	CurrentTerm int
	VotedFor    string
	chan_commits chan int  //更新了commit后，激活这个管道
	chan_timeout chan bool //定时管道
}

func (this *Follower) Init(role_chan chan RoleState) error {
	this.chan_role = role_chan

	this.active = new(utils.AtomicBool)
	this.active.Set()

	this.chan_commits = make(chan int)
	this.chan_timeout = make(chan bool)
	return nil
}

func (this *Follower) StartAllService() {
	go this.startLogApplService() //日志应用服务
	go this.startTimeOutService() //计时服务
}

//定时服务，超时即切换到candidate状态
func (this *Follower) startTimeOutService() {
	for {
		if !this.active.IsSet() {
			break
		}
		ot := time.Duration(rand.Intn(settings.TIMEOUT_MAX-settings.TIMEOUT_MIN) + settings.TIMEOUT_MIN)
		select {
		case <-this.chan_timeout:
			//do nothing
		case <-time.After(ot * time.Millisecond):
			if !this.active.IsSet() {
				rolestate := RoleState{settings.CANDIDATE, this.CurrentTerm}
				this.chan_role <- rolestate
			}
			break
		}
	}
}

//日志应用服务(更新lastApplied)
func (this *Follower) startLogApplService() {
	for {
		if !this.active.IsSet() {
			return
		}
		for this.LastApplied < this.CommitIndex {
			_log := this.Logs.Get(this.LastApplied + 1)
			db.WriteToDisk(_log.Command)
			this.LastApplied++
		}

		select {
		case <-this.chan_commits:
			//do nothing
		case <-time.After(time.Millisecond * time.Duration(settings.NEWLOG_WAIT)):
			//do nothing
		}
	}
}

func (this *Follower) HandleVoteReq(args0 VoteReqArg, args1 *VoteAckArg) error {
	if !this.active.IsSet() {
		return nil
	}

	//过期leader直接拒绝
	if this.CurrentTerm > args0.Term {
		args1.Term = this.CurrentTerm
		args1.VoteGranted = false
		return nil
	}

	this.CurrentTerm = args0.Term

	voteGranted := false
	if this.VotedFor == args0.CandidateID || this.VotedFor == "" { //比较谁包含的日志记录更新更长
		log_length := this.Logs.Size()
		last_log := this.Logs.Get(log_length - 1)

		switch t := last_log.Term - args0.LastLogTerm {
		case t < 0:
			voteGranted = true
		case t > 0:
			voteGranted = false
		case t == 0:
			if log_length > args0.LastLogIndex+1 {
				voteGranted = false
			} else {
				voteGranted = true
			}
		}
	}
	if voteGranted {
		this.VotedFor = args0.CandidateID
		args1.Term = this.CurrentTerm
		args1.VoteGranted = true
	} else {
		args1.Term = this.CurrentTerm
		args1.VoteGranted = false
	}

	return nil
}

func (this *Follower) HandleAppendLogReq(args0 LogAppArg, args1 *LogAckArg) error {
	if !this.active.IsSet() {
		return nil
	}

	for {
		if !this.active.IsSet() {
			return nil
		}
		select {
		case this.chan_timeout <- true: //重置定时器
		case time.After(time.Millisecond * time.Duration(settings.CHANN_WAIT)):
			continue
		}
	}

	//收到了过期leader的log_rpc
	if args0.Term > this.CurrentTerm {
		args1.Term = this.CurrentTerm
		args1.Success = false
		return nil
	}

	this.VotedFor = args0.LeaderID

	this.CurrentTerm = args0.Term

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

func (this *Follower) handleHeartBeat(args0 LogAppArg, args1 *LogAckArg) error {
	if !this.active.IsSet() {
		return nil
	}

	if args0.LeaderCommitIdx > this.CommitIndex {
		log_size := this.Logs.Size()
		if args0.LeaderCommitIdx < log_size {
			this.CommitIndex = args0.LeaderCommitIdx
		} else {
			this.CommitIndex = log_size
		}

		go func() {
			for {
				if !this.active.IsSet() {
					return
				}
				select {
				case this.chan_commits <- this.CommitIndex:
					break
				case time.After(time.Millisecond * time.Duration(settings.COMMIT_WAIT)):
					continue
				}
			}

		}()

	}
	return nil
}

func (this *Follower) HandleCommandReq(cmd string, ok *bool, leaderIP *string) error {
	*ok = false
	*leaderIP = this.VotedFor
	return nil
}

package servers

import (
	"SimpleRaft/db"
	"time"
	"SimpleRaft/settings"
	"math/rand"
	"SimpleRaft/utils"
	"log"
)

type Follower struct {
	*BaseRole

	//persistent state
	CurrentTerm int
	VotedFor    string

	chan_timeout chan bool //定时管道

	//角色是否处于激活状态，供角色启动的子协程参考
	//用指针防止copy
	active    *utils.AtomicBool
	chan_role *chan RoleState //角色管道
}

func (this *Follower) Init() error {
	tmp := make(chan RoleState)
	this.chan_role = &tmp

	this.active = new(utils.AtomicBool)
	this.active.Set()

	this.chan_timeout = make(chan bool)

	log.Printf("FOLLOWER(%d)：初始化...\n", this.CurrentTerm)

	return nil
}

func (this *Follower) SetAlive(alive bool) {
	this.active.SetTo(alive)
}

func (this *Follower) GetAlive() bool {
	return this.active.IsSet()
}

func (this *Follower) GetRoleChan() *chan RoleState {
	return this.chan_role
}

func (this *Follower) StartAllService() {
	go this.startTimeOutService() //计时服务
	go this.startLogApplService() //日志应用服务
}

//定时服务，超时即切换到candidate状态
func (this *Follower) startTimeOutService() {
	for {
		if !this.active.IsSet() {
			break
		}
		ot := time.Duration(rand.Intn(settings.TIMEOUT_MAX-settings.TIMEOUT_MIN) + settings.TIMEOUT_MIN)
		log.Printf("FOLLOWER(%d)：启动计时服务(%d毫秒)...\n", this.CurrentTerm, ot)
		select {
		case st := <-this.chan_timeout:
			if st {
				log.Printf("FOLLOWER(%d)：暂停计时器...\n", this.CurrentTerm)
				return
			} else {
				log.Printf("FOLLOWER(%d)：重置计时器...\n", this.CurrentTerm)
				break
			}
		case <-time.After(ot * time.Millisecond):
			if this.active.IsSet() {
				log.Printf("FOLLOWER(%d)：Leader(%s)一直未响应，开始选举...\n", this.CurrentTerm, this.VotedFor)
				rolestate := RoleState{settings.CANDIDATE, this.CurrentTerm}
				*this.chan_role <- rolestate
			} else {
				log.Printf("FOLLOWER(%d)：计时服务终止！！！\n", this.CurrentTerm)
				return
			}
		}
	}
}

//日志应用服务(更新lastApplied)
func (this *Follower) startLogApplService() {
	log.Printf("FOLLOWER(%d)：启动日志应用服务...\n", this.CurrentTerm)
	for {
		if !this.active.IsSet() {
			break
		}
		for this.LastApplied < this.CommitIndex {
			this.LastApplied++
			_log := this.Logs.Get(this.LastApplied)
			db.WriteToDisk(_log.Command)
			log.Printf("FOLLOWER(%d): <NUM: %d, TERM: %d, CMD: %s> 已写到日志文件\n",
				this.CurrentTerm, this.LastApplied, _log.Term, _log.Command)
		}

		time.Sleep(time.Millisecond * time.Duration(settings.COMMIT_WAIT))
	}
	log.Printf("FOLLOWER(%d)：日志应用服务终止！！！\n", this.CurrentTerm)
}

func (this *Follower) HandleVoteReq(args0 VoteReqArg, args1 *VoteAckArg) error {
	if !this.active.IsSet() {
		return nil
	}

	log.Printf("FOLLOWER(%d)：收到%s的投票请求...\n", this.CurrentTerm, args0.CandidateID)

	//过期leader直接拒绝
	if this.CurrentTerm > args0.Term {
		args1.Term = this.CurrentTerm
		args1.VoteGranted = false
		log.Printf("FOLLOWER(%d)：%s为过期的CANDIDATE，拒绝为其投票！！！\n", this.CurrentTerm, args0.CandidateID)
		return nil
	}

	this.CurrentTerm = args0.Term

	voteGranted := false
	if this.VotedFor == args0.CandidateID || this.VotedFor == "" { //比较谁包含的日志记录更新更长
		lastLogIndex := this.Logs.Size()
		lastLog := this.Logs.Get(lastLogIndex)

		t := lastLog.Term - args0.LastLogTerm
		switch {
		case t < 0:
			voteGranted = true
		case t > 0:
			voteGranted = false
			log.Printf("FOLLOWER(%d)：不支持%s(Term:%d过期)\n", this.CurrentTerm, args0.CandidateID, args0.LastLogTerm)
		case t == 0:
			if lastLogIndex > args0.LastLogIndex {
				voteGranted = false
				log.Printf("FOLLOWER(%d)：不支持%s(日志不够新)\n", this.CurrentTerm, args0.CandidateID)
			} else {
				voteGranted = true
			}
		}
	}
	if voteGranted {
		this.VotedFor = args0.CandidateID
		args1.Term = this.CurrentTerm
		args1.VoteGranted = true
		log.Printf("FOLLOWER(%d)：投票给%s\n", this.CurrentTerm, args0.CandidateID)
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
		case this.chan_timeout <- true:
			defer func() {
				go this.startTimeOutService()
			}()
			goto EndFor
		case <-time.After(time.Millisecond * time.Duration(settings.CHANN_WAIT)):
			break //跳出select
		}
	}
EndFor:

//收到了过期leader的log_rpc
	if args0.Term < this.CurrentTerm {
		args1.Term = this.CurrentTerm
		args1.Success = false
		log.Printf("FOLLOWER(%d)：收到过期leader(%s)的append_log请求\n", this.CurrentTerm, args0.LeaderID)
		return nil
	}

	this.VotedFor = args0.LeaderID

	this.CurrentTerm = args0.Term

	//来自leader的心跳信息(更新commit index)
	if args0.Entries == nil || len(args0.Entries) == 0 {
		this.handleHeartBeat(args0)
	} else {
		if args0.PreLogIndex > 0 {
			preLog := this.Logs.Get(args0.PreLogIndex)
			if preLog.Term != args0.PreLogTerm || preLog.Command == "" {
				if preLog.Command == "" {
					log.Printf("FOLLOWER(%d)：日志落后于leader(%s)：(FOLLOWER:%d/LEADER:%d)\n",
						this.CurrentTerm, args0.LeaderID, this.Logs.Size(), args0.PreLogIndex)
				}

				if preLog.Term != 0 && preLog.Term != args0.PreLogTerm {
					log.Printf("FOLLOWER(%d)：与leader(%s)的日志信息不一致(%d/%d)\n",
						this.CurrentTerm, args0.LeaderID, preLog.Term, args0.PreLogTerm)
					this.Logs.RemoveFrom(args0.PreLogIndex)
				}

				args1.LastLogIndex = this.Logs.Size()
				args1.Term = this.CurrentTerm
				args1.Success = false
				return nil
			}
		}
	}

	lastLogIdx := this.Logs.Extend(args0.PreLogIndex, args0.Entries)
	args1.LastLogIndex = lastLogIdx
	args1.Term = this.CurrentTerm
	args1.Success = true

	time.Sleep(time.Millisecond * time.Duration(rand.Intn(10000)))

	if args0.Entries != nil && len(args0.Entries) > 0 {
		log.Printf("FOLLOWER(%d)：收到leader(%s)的新日志(Totals: %d, PreTerm: %d, PreIndex: %d, Size: %d)\n",
			this.CurrentTerm, args0.LeaderID, this.Logs.Size(), args0.PreLogTerm, args0.PreLogIndex, len(args0.Entries))
	} else {
		log.Printf("FOLLOWER(%d)：收到leader(%s)的心跳信息\n", this.CurrentTerm, args0.LeaderID)
	}

	return nil
}

func (this *Follower) handleHeartBeat(args0 LogAppArg) error {
	if !this.active.IsSet() {
		return nil
	}

	if args0.LeaderCommitIdx > this.CommitIndex {
		log_size := this.Logs.Size()
		if args0.LeaderCommitIdx > log_size {
			this.CommitIndex = log_size
		} else {
			this.CommitIndex = args0.LeaderCommitIdx
		}
	}
	return nil
}

func (this *Follower) HandleCommandReq(cmd string, cmdAck *CommandAck) error {
	cmdAck.Ok = false
	cmdAck.LeaderIP = this.VotedFor
	return nil
}

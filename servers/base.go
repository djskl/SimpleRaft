package servers

import (
	"SimpleRaft/utils"
	"SimpleRaft/clog"
)

type RaftServer interface {
	Init() error												//初始化操作
	HandleVoteReq(args0 VoteReqArg, args1 *VoteAckArg) error 	//处理RequestVote RPC
	HandleAppendLogReq(args0 LogAppArg, args1 *LogAckArg) error	//处理AppendEntries RPC
	HandleCommandReq(cmds string, ok *bool) error				//处理用户(client)请求
	SetAlive(alive bool)										//改变角色的存活状态
	StartAllService()											//启动当前角色下的所有服务
}

type BaseRole struct {
	IP string

	//persistent state
	CurrentTerm int
	VotedFor    string

	Logs *clog.Manager

	//volatile state
	CommitIndex int
	LastApplied int

	//角色是否处于激活状态，供角色启动的子协程参考
	active *utils.AtomicBool	//用指针防止copy

	chan_role chan int	//角色管道，由manager初始化，所有角色共享一个
}

func (this *BaseRole) init(chan_role chan int) {
	if chan_role == nil {
		panic("Role Chan has not been initialized")
	}
	this.chan_role = chan_role
	this.Logs = new(clog.Manager)
	this.Logs.Init()
	this.active.Set()
}

func (this *BaseRole) GetAlive() bool {
	return this.active.IsSet()
}

func (this *BaseRole) SetAlive(alive bool) {
	this.active.UnSet()
}

type VoteReqArg struct {
	Term int
	CandidateID string
	LastLogIndex int
	LastLogTerm int
}

type VoteAckArg struct {
	Term int
	VoteGranted bool
}

type LogAppArg struct {
	Term int
	LeaderID string
	PreLogIndex int
	PreLogTerm int
	LeaderCommitIdx int
	Entries []clog.Item
}

type LogAckArg struct {
	Term int
	Success bool
	LastLogIndex int
}


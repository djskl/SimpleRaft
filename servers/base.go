package servers

import (
	"SimpleRaft/utils"
	"SimpleRaft/clog"
)

type RaftServer interface {
	Init(role_chan chan RoleState) error                                  //初始化操作										//改变角色的存活状态
	StartAllService()                                               //启动当前角色下的所有服务
	HandleVoteReq(args0 VoteReqArg, args1 *VoteAckArg) error        //处理RequestVote RPC
	HandleAppendLogReq(args0 LogAppArg, args1 *LogAckArg) error     //处理AppendEntries RPC
	HandleCommandReq(cmds string, ok *bool, leaderIP *string) error //处理用户(client)请求
}

type BaseRole struct {
	IP string

	Logs *clog.Manager

	//volatile state
	CommitIndex int
	LastApplied int

	//角色是否处于激活状态，供角色启动的子协程参考
	//用指针防止copy
	active *utils.AtomicBool

	chan_role chan RoleState
}

func (this *BaseRole) Init() {
	this.Logs = new(clog.Manager)
	this.Logs.Init()
}

func (this *BaseRole) GetAlive() bool {
	return this.active.IsSet()
}

func (this *BaseRole) SetAlive(alive bool) {
	this.active.UnSet()
}

type VoteReqArg struct {
	Term         int
	CandidateID  string
	LastLogIndex int
	LastLogTerm  int
}

type VoteAckArg struct {
	Term        int
	VoteGranted bool
}

type LogAppArg struct {
	Term            int
	LeaderID        string
	PreLogIndex     int
	PreLogTerm      int
	LeaderCommitIdx int
	Entries         []clog.Item
}

type LogAckArg struct {
	Term         int
	Success      bool
	LastLogIndex int
}

type RoleState struct {
	Role int
	Term int
}

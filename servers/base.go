package servers

type RaftServer interface {
	Init() error												//初始化操作
	HandleVoteReq(args0 VoteReqArg, args1 *VoteAckArg) error 	//处理RequestVote RPC
	HandleAppendLogReq(args0 LogAppArg, args1 *LogAckArg) error	//处理AppendEntries RPC
	HandleCommandReq(cmds string, ok *bool) error				//处理用户(client)请求
	SetAlive(alive bool)										//改变角色的存活状态
}

type BaseRole struct {
	IP string
	//persistent state
	CurrentTerm int
	VotedFor    string
	Logs        []LogItem
	//volatile state
	CommitIndex int
	LastApplied int
	//角色是否处于激活状态，供角色启动的子协程参考
	IsAlive bool

}

//日志项
type LogItem struct {
	Term    int
	Command string
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
	Entries []LogItem
}

type LogAckArg struct {
	Term int
	Success bool
}


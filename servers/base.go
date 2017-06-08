package servers

type RaftServer interface {
	Init() error												//初始化操作
	HandleVoteReq(args0 VoteReqArg, args1 *VoteAckArg) error 	//处理RequestVote RPC
	HandleAppendLogReq(args0 LogAppArg, args1 *LogAckArg) error	//处理AppendEntries RPC
	HandleCommandReq(cmds string, ok *bool) error				//处理用户(client)请求
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


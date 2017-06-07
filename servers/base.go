package servers

type VoteReqArg struct {

}

type VoteAckArg struct {

}

type LogAppArg struct {

}

type LogAckArg struct {

}

type RaftServer interface {
	Init() error												//初始化操作
	HandleVoteReq(args0 VoteReqArg, args1 *VoteAckArg) error 	//处理RequestVote RPC
	HandleAppendLogReq(args0 LogAppArg, args1 *LogAckArg) error	//处理AppendEntries RPC
	HandleCommandReq(cmds string, ok *bool) error				//处理用户(client)请求
}

//日志项
type LogItem struct {
	Term    int
	Command string
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


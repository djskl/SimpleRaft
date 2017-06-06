package servers

type RaftServer interface {
	Init()				//初始化操作
	HandleVoteReq()		//处理RequestVote RPC
	HandleReceivedMsg()	//处理AppendEntries RPC
	HandleMsgAck()		//处理针对AppendEntries RPC的确认信息
}

//日志项
type LogItem struct {
	Term    int
	Command string
}

type BaseRole struct {
	ID string
	//persistent state
	CurrentTerm int
	VotedFor    string
	Logs        []LogItem
	//volatile state
	CommitIndex int
	LastApplied int
}


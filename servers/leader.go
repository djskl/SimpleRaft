package servers

type Leader struct {
	*BaseRole
	NextIndex  map[string]int
	MatchIndex map[string]int
}

func (this *Leader) Init() error {

	return nil
}

func (this *Leader) HandleVoteReq(args0 VoteReqArg, args1 *VoteAckArg) error {
	return nil
}

func (this *Leader) HandleAppendLogReq(args0 LogAppArg, args1 *LogAckArg) error {
	return nil
}

func (this *Leader) HandleCommandReq(cmds string, ok *bool) error {
	return nil
}
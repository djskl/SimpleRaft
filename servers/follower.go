package servers

type Follower struct {
	*BaseRole
}

func (this *Follower) Init(role_chan chan int) error {
	return nil
}

func (this *Follower) StartAllService() {

}

func (this *Follower) HandleVoteReq(args0 VoteReqArg, args1 *VoteAckArg) error {
	return nil
}

func (this *Follower) HandleAppendLogReq(args0 LogAppArg, args1 *LogAckArg) error {
	return nil
}

func (this *Follower) HandleCommandReq(cmd string, ok *bool) error {
	return nil
}

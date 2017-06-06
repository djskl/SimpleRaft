package servers

type Leader struct {
	*BaseRole
	NextIndex  []int
	MatchIndex []int
}


func (this *Leader) Init() {

}

func (this *Leader) HandleVoteReq() {

}

func (this *Leader) HandleReceivedMsg() {

}

func (this *Leader) HandleMsgAck() {

}
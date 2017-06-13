package servers

import (
	"net/rpc"
	"SimpleRaft/settings"
	"log"
	"sync/atomic"
)

type Candidate struct {
	*BaseRole
	chan_timeout chan bool
	total_votes  *int32
}

func (this *Candidate) Init(role_chan chan int) error {
	this.BaseRole.init(role_chan)
	this.chan_timeout = make(chan bool)
	this.total_votes = new(int32)
	atomic.AddInt32(this.total_votes, 1)	//先投自己一票
	return nil
}

func (this *Candidate) requestVote(ip string, logIdx int, logTerm int) {

	voteReq := VoteReqArg{
		this.CurrentTerm,
		this.IP,
		logIdx,
		logTerm,
	}

	voteAck := new(VoteAckArg)

	var client *rpc.Client
	var err error
	for {
		if !this.active.IsSet() {
			return
		}
		client, err = rpc.DialHTTP("tcp", ip+":"+settings.SERVERPORT)
		if err != nil {
			log.Printf("failed to connect %s from %s\n", ip, this.IP)
			continue
		}
		break
	}

	for {
		if !this.active.IsSet() {
			client.Close()
			return
		}
		err = client.Call("RaftManager.Vote", voteReq, voteAck)
		if err != nil {
			log.Printf("appendLog failed: %s--->%s", this.IP, ip)
			continue
		}
		break
	}
	client.Close()
}

func (this *Candidate) handleVoteAck(voteAck *VoteAckArg) {
	//收到了比自己大的term直接转换为follower
	if voteAck.Term > this.CurrentTerm {
		this.CurrentTerm = voteAck.Term
		this.chan_role <- settings.FOLLOWER
		return
	}

	// 如果因为日志不够up-to-date而失败，
	// 此时仍有可能被选为leader，不能直接降为follower
	if voteAck.VoteGranted {
		atomic.AddInt32(this.total_votes, 1)
		if(*this.total_votes > 2){	//选举成功了
			this.chan_role <- settings.LEADER
			return
		}
	}
}

func (this *Candidate) HandleVoteReq(args0 VoteReqArg, args1 *VoteAckArg) error {

	args1.Term = this.CurrentTerm
	args1.VoteGranted = false	//candidate只有被选举权，没有选举权

	//收到了比自己大的term直接转换为follower
	if args0.Term > this.CurrentTerm {
		this.CurrentTerm = args0.Term
		this.chan_role <- settings.FOLLOWER
	}

	return nil
}





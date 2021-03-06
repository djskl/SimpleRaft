package servers

import (
	"net/rpc"
	"SimpleRaft/settings"
	"log"
	"sync/atomic"
	"time"
	"math/rand"
	"SimpleRaft/utils"
)

type Candidate struct {
	*BaseRole

	CurrentTerm int
	VotedFor    string

	chan_timeout chan bool
	total_votes  *int32

	AllServers []string

	//角色是否处于激活状态，供角色启动的子协程参考
	//用指针防止copy
	active    *utils.AtomicBool
	chan_role *chan RoleState //角色管道
}

func (this *Candidate) Init() error {
	tmp := make(chan RoleState)
	this.chan_role = &tmp

	this.VotedFor = ""

	this.active = new(utils.AtomicBool)
	this.active.Set()

	this.chan_timeout = make(chan bool)

	this.CurrentTerm += 1
	this.total_votes = new(int32)
	atomic.AddInt32(this.total_votes, 1) //先投自己一票

	log.Printf("CANDIDATE(%d)：初始化...\n", this.CurrentTerm)

	return nil
}

func (this *Candidate) SetAlive(alive bool) {
	this.active.SetTo(alive)
}

func (this *Candidate) GetAlive() bool {
	return this.active.IsSet()
}

func (this *Candidate) GetRoleChan() *chan RoleState {
	return this.chan_role
}

func (this *Candidate) StartAllService() {
	go this.startVoteService()
}

func (this *Candidate) startVoteService() {

	log.Printf("CANDIDATE(%d)：启动投票服务...\n", this.CurrentTerm)

	lastLogIndex := this.Logs.Size()
	lastLog, err := this.Logs.Get(lastLogIndex)

	var lastLogTerm int

	if err != nil {
		lastLogTerm = 0
	} else {
		lastLogTerm = lastLog.Term
	}

	for idx := 0; idx < len(this.AllServers); idx++ {
		ip := this.AllServers[idx]
		if ip == this.IP {
			continue
		}
		go this.requestVote(ip, lastLogIndex, lastLogTerm)
	}

	ot := time.Duration(rand.Intn(settings.TIMEOUT_MAX-settings.TIMEOUT_MIN) + settings.TIMEOUT_MIN)

	<-time.After(ot * time.Millisecond)

	if !this.active.IsSet() {
		return
	}

	log.Printf("CANDIDATE(%d)：未选出Leader，重新选举...\n", this.CurrentTerm)
	rolestate := RoleState{settings.CANDIDATE, this.CurrentTerm}
	*this.chan_role <- rolestate
}

func (this *Candidate) requestVote(ip string, logIdx int, logTerm int) {
	if !this.active.IsSet() {
		return
	}

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
			log.Printf("CANDIDATE(%d)：无法与%s建立连接！！！\n", this.CurrentTerm, ip)
			time.Sleep(time.Second * time.Duration(settings.RPC_WAIT))
			continue
		}
		break
	}

	for {
		if !this.active.IsSet() {
			client.Close()
			return
		}
		err = client.Call("RaftRPC.Vote", voteReq, voteAck)
		if err != nil {
			log.Printf("CANDIDATE(%d)：调用%s的Vote方法失败！！！\n", this.CurrentTerm, ip)
			time.Sleep(time.Second * time.Duration(settings.RPC_WAIT))
			continue
		}
		break
	}
	client.Close()

	log.Printf("CANDIDATE(%d)：向%s发送投票请求...\n", this.CurrentTerm, ip)

	this.handleVoteAck(voteAck)

}

func (this *Candidate) handleVoteAck(voteAck *VoteAckArg) {
	if !this.active.IsSet() {
		return
	}

	//收到了比自己大的term直接转换为follower
	if voteAck.Term > this.CurrentTerm {
		log.Printf("CANDIDATE(%d)：Term过期了，转为Follower...\n", this.CurrentTerm)
		this.CurrentTerm = voteAck.Term
		rolestate := RoleState{settings.FOLLOWER, this.CurrentTerm}
		if !this.active.IsSet() {
			return
		}
		select {
		case *this.chan_role <- rolestate:
			//do nothing
		case <-time.After(time.Millisecond * time.Duration(settings.CHANN_WAIT)):
			//
		}
		return
	}

	// 如果因为日志不够up-to-date而失败，
	// 此时仍有可能被选为leader，不能直接降为follower
	if voteAck.VoteGranted {
		log.Printf("CANDIDATE(%d)：获得一票，当前票数：%d\n", this.CurrentTerm, *this.total_votes)
		atomic.AddInt32(this.total_votes, 1)
		if *this.total_votes >= settings.MAJORITY { //选举成功了
			log.Printf("CANDIDATE(%d)：选举成功(票数:%d)，转为Leader\n", this.CurrentTerm, *this.total_votes)
			rolestate := RoleState{settings.LEADER, this.CurrentTerm}
			if !this.active.IsSet() {
				return
			}
			select {
			case *this.chan_role <- rolestate:
				//
			case <-time.After(time.Millisecond * time.Duration(settings.CHANN_WAIT)):
				//
			}
		}
	}
}

func (this *Candidate) HandleVoteReq(args0 VoteReqArg, args1 *VoteAckArg) error {

	log.Printf("CANDIDATE(%d)：收到投票请求...\n", this.CurrentTerm)

	args1.Term = this.CurrentTerm
	args1.VoteGranted = false //candidate只有被选举权，没有选举权

	//收到了比自己大的term直接转换为follower
	if args0.Term > this.CurrentTerm {
		log.Printf("CANDIDATE(%d)：Term过期了，转为Follower...\n", this.CurrentTerm)
		this.CurrentTerm = args0.Term
		rolestate := RoleState{settings.FOLLOWER, this.CurrentTerm}
		if !this.active.IsSet() {
			return nil
		}
		select {
		case *this.chan_role <- rolestate:
			//do nothing
		case <-time.After(time.Millisecond * time.Duration(settings.CHANN_WAIT)):
			//
		}
	}
	return nil
}

//处于Candidate状态的机器暂时不可用
func (this *Candidate) HandleCommandReq(cmds string, cmdAck *CommandAck) error {
	cmdAck.Ok = false
	cmdAck.LeaderIP = this.VotedFor
	return nil
}

func (this *Candidate) HandleAppendLogReq(args0 LogAppArg, args1 *LogAckArg) error {
	log.Printf("CANDIDATE(%d)：收到AppendLog请求...\n", this.CurrentTerm)
	args1.Term = this.CurrentTerm
	args1.Success = false

	if args0.Term >= this.CurrentTerm {
		log.Printf("CANDIDATE(%d)：Term过期了，转为Follower...\n", this.CurrentTerm)
		this.CurrentTerm = args0.Term
		rolestate := RoleState{settings.FOLLOWER, this.CurrentTerm}
		if !this.active.IsSet() {
			return nil
		}
		select {
		case *this.chan_role <- rolestate:
			//do nothing
		case <-time.After(time.Millisecond * time.Duration(settings.CHANN_WAIT)):
			//do nothing
		}
	}
	return nil
}

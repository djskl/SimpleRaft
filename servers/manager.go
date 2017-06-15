package servers

import (
	"SimpleRaft/settings"
	"log"
	"SimpleRaft/db"
)

type RaftManager struct {
	br *BaseRole
	rs RaftServer

	AllServers []string
	CurrentIP  string

	role_stop chan bool
	role_over chan bool

}

func (this *RaftManager) Init() {
	this.role_stop = make(chan bool)
	this.role_over = make(chan bool)

	db.Init()

	this.AllServers = db.LoadServerIPS(settings.IPFILE)
	this.CurrentIP = this.AllServers[0]

	this.br = &BaseRole{IP: this.CurrentIP}
	this.br.Init()

	this.rs = &Follower{BaseRole: this.br}
	this.rs.Init()
	this.rs.StartAllService()
	this.StartRoleService()
}

func (this *RaftManager) StartRoleService() {
	role_chan := this.rs.GetRoleChan()
	go func() {
		log.Println("MANAGER：启动角色监听服务...")
		for {
			select {
			case <- this.role_stop:
				this.role_over <- true
				return
			case role := <-role_chan:
				go this.ConvertToRole(role)
			}
		}
	}()
}

func (this *RaftManager) RestartRoleService() {
	role_chan := this.rs.GetRoleChan()
	go func() {
		this.role_stop <- true
		<- this.role_over
		log.Println("MANAGER：重启角色监听服务...")
		for {
			select {
			case <- this.role_stop:
				this.role_over <- true
				return
			case role := <-role_chan:
				this.ConvertToRole(role)
			}
		}
	}()
}

func (this *RaftManager) ConvertToRole(state RoleState) {
	this.rs.SetAlive(false) //注销当前角色
	switch state.Role {
	case settings.LEADER:
		this.rs = &Leader{
			BaseRole:    this.br,
			CurrentTerm: state.Term,
			AllServers:  this.AllServers,
		}
	case settings.CANDIDATE:
		this.rs = &Candidate{
			BaseRole:    this.br,
			CurrentTerm: state.Term,
			AllServers:  this.AllServers,
		}
	case settings.FOLLOWER:
		this.rs = &Follower{
			BaseRole:    this.br,
			CurrentTerm: state.Term,
		}
	default:
		log.Fatalf("MANAGER：未知角色：%d\n", state.Role)
	}
	this.rs.Init()
	this.rs.StartAllService()
	this.RestartRoleService()
}

//Vote is used to respond to RequestVote RPC
func (this *RaftManager) Vote(args0 VoteReqArg, args1 *VoteAckArg) error {
	err := this.rs.HandleVoteReq(args0, args1)
	return err
}

//AppendLog is used to respond to AppendEntries RPC（with filled log）
func (this *RaftManager) AppendLog(args0 LogAppArg, args1 *LogAckArg) error {
	err := this.rs.HandleAppendLogReq(args0, args1)
	return err
}

//Command is used to interact with users
func (this *RaftManager) Command(cmds string, cmdAck *CommandAck) error {
	err := this.rs.HandleCommandReq(cmds, &cmdAck.ok, &cmdAck.leaderIP)
	return err
}

type RaftRPC struct {
	RM *RaftManager
}

//Vote is used to respond to RequestVote RPC
func (this *RaftRPC) Vote(args0 VoteReqArg, args1 *VoteAckArg) error {
	err := this.RM.Vote(args0, args1)
	return err
}

//AppendLog is used to respond to AppendEntries RPC（with filled log）
func (this *RaftRPC) AppendLog(args0 LogAppArg, args1 *LogAckArg) error {
	err := this.RM.AppendLog(args0, args1)
	return err
}

//Command is used to interact with users
func (this *RaftRPC) Command(cmds string, cmdAck *CommandAck) error {
	err := this.RM.Command(cmds, cmdAck)
	return err
}

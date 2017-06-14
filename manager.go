package main

import (
	"SimpleRaft/servers"
	"SimpleRaft/settings"
	"log"
)

type RaftManager struct {
	br *servers.BaseRole
	rs servers.RaftServer

	chan_role chan servers.RoleState //角色管道
}

func (this *RaftManager) init() {
	this.br = &servers.BaseRole{IP: settings.CurrentIP}
	this.br.Init()
	this.rs = &servers.Follower{BaseRole: this.br}
	this.chan_role = make(chan servers.RoleState)
	this.rs.Init(this.chan_role)
}

//通过管道chan_role监听角色变化事件
func (this *RaftManager) StartRoleService() {
	go func() {
		for {
			role := <-this.chan_role
			this.convertToRole(role)
		}
	}()
}

func (this *RaftManager) convertToRole(state servers.RoleState) {
	this.br.SetAlive(false) //注销当前角色
	switch state.Role {
	case settings.LEADER:
		this.rs = &servers.Leader{BaseRole: this.br, CurrentTerm: state.Term}
	case settings.FOLLOWER:
		this.rs = &servers.Follower{BaseRole: this.br, CurrentTerm: state.Term}
	case settings.CANDIDATE:
		this.rs = &servers.Candidate{BaseRole: this.br, CurrentTerm: state.Term}
	default:
		log.Fatalf("the role: %d doesn't exist", state.Role)
	}
	this.rs.Init(this.chan_role)
}

//Vote is used to respond to RequestVote RPC
func (this *RaftManager) Vote(args0 servers.VoteReqArg, args1 *servers.VoteAckArg) error {
	err := this.rs.HandleVoteReq(args0, args1)
	return err
}

//AppendLog is used to respond to AppendEntries RPC（with filled log）
func (this *RaftManager) AppendLog(args0 servers.LogAppArg, args1 *servers.LogAckArg) error {
	err := this.rs.HandleAppendLogReq(args0, args1)
	return err
}

//Command is used to interact with users
func (this *RaftManager) Command(cmds string, ok *bool, leaderIP *string) error {
	err := this.rs.HandleCommandReq(cmds, ok, leaderIP)
	return err
}

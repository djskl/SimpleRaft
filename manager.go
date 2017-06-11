package main

import (
	"SimpleRaft/servers"
	"SimpleRaft/utils"
	"SimpleRaft/settings"
	"log"
)

type RaftManager struct {
	br *servers.BaseRole
	rs servers.RaftServer

	chan_role chan int	//角色管道
}

func (this *RaftManager) init() {
	this.br = &servers.BaseRole{IP: settings.CurrentIP}
	this.rs = &servers.Follower{this.br}
	this.chan_role = make(chan int)
	this.rs.Init(this.chan_role)
}

//通过管道chan_role监听角色变化事件
func (this *RaftManager) StartRoleService()  {
	for {
		role := <- this.chan_role
		this.convertToRole(role)
	}
}

func (this *RaftManager) convertToRole(role int) {
	this.br.SetAlive(false) //注销当前角色
	switch role {
	case settings.LEADER:
		this.rs = &servers.Leader{BaseRole: this.br}
	case settings.FOLLOWER:
		this.rs = &servers.Follower{BaseRole: this.br}
	case settings.CANDIDATE:
		//TODO
	default:
		log.Fatalf("the role: %d doesn't exist", role)
	}
	this.rs.Init(this.chan_role)
}

//Vote is used to respond to RequestVote RPC
func (this *RaftManager) Vote() error {
	return nil
}

//AppendLog is used to respond to AppendEntries RPC（with filled log）
func (this *RaftManager) AppendLog() error {
	return nil
}

//HeartBeat is used to respond to AppendEntries RPC（with empty log）
func (this *RaftManager) HeartBeat() error {
	return nil
}

//Command is used to interact with users
func (this *RaftManager) Command() error {
	return nil
}

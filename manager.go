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
	server_id := utils.GetUUID(settings.UUIDSIZE)
	this.br = &servers.BaseRole{IP: server_id}
	this.rs = &servers.Follower{this.br}
	this.rs.Init()
}

//通过管道chan_role监听角色变化事件
func (this *RaftManager) StartRoleService()  {
	for {
		role := <- this.chan_role
		this.convertToRole(role)
	}
}

func (this *RaftManager) convertToRole(role int) {
	this.rs.SetAlive(false) //注销当前角色
	switch role {
	case settings.LEADER:
		this.rs = &servers.Leader{BaseRole: this.br}
	case settings.FOLLOWER:
		//TODO
	case settings.CANDIDATE:
		//TODO
	default:
		log.Fatalf("the role: %d doesn't exist", role)
	}
	this.rs.Init()
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

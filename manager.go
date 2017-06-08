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
}

func (m *RaftManager) init() {
	server_id := utils.GetUUID(settings.UUIDSIZE)
	m.br = &servers.BaseRole{IP: server_id}
	m.rs = servers.Follower{m.br}
	m.rs.Init()
}

func (m *RaftManager) convertToRole(role int) {
	switch role {
	case settings.LEADER:
		m.rs = servers.Leader{BaseRole: m.br}
	case settings.FOLLOWER:
		m.rs = servers.Follower{BaseRole: m.br}
	case settings.CANDIDATE:
		//TODO
	default:
		log.Fatalf("the role: %d doesn't exist", role)
	}
	m.rs.Init()
}

//Vote is used to respond to RequestVote RPC
func (m *RaftManager) Vote() error {
	return nil
}

//AppendLog is used to respond to AppendEntries RPC（with filled log）
func (m *RaftManager) AppendLog() error {
	return nil
}

//HeartBeat is used to respond to AppendEntries RPC（with empty log）
func (m *RaftManager) HeartBeat() error {
	return nil
}

//Command is used to interact with users
func (m *RaftManager) Command() error {
	return nil
}


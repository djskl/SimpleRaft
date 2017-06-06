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
	m.br = &servers.BaseRole{ID: server_id}
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



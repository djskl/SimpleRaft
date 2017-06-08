package servers

import (
	"SimpleRaft/settings"
	"log"
	"fmt"
	"net/rpc"
)

type Leader struct {
	*BaseRole
	NextIndex  map[string]int
	MatchIndex map[string]int
}

func (this *Leader) Init() error {
	return nil
}

func (this *Leader) HandleVoteReq(args0 VoteReqArg, args1 *VoteAckArg) error {
	return nil
}

func (this *Leader) HandleAppendLogReq(args0 LogAppArg, args1 *LogAckArg) error {
	return nil
}

func (this *Leader) HandleCommandReq(cmds string, ok *bool) error {
	for idx := 0; idx < len(settings.AllServers); idx++ {
		ip := settings.AllServers[idx]
		if ip == this.IP {
			continue
		}
	}
	return nil
}

func (this *Leader) replicateLog(ip string) {

	nextLogIdx := this.NextIndex[ip]
	preLog := this.Logs[nextLogIdx-1]
	toSendEntries := this.Logs[nextLogIdx:]

	inArg := LogAppArg{
		Term: this.CurrentTerm,
		LeaderID: this.IP,
		LeaderCommitIdx: this.CommitIndex,
		PreLogIndex: nextLogIdx - 1,
		PreLogTerm: preLog.Term,
		Entries: toSendEntries,
	}
	client, err := rpc.DialHTTP("tcp", ip+":"+settings.SERVERPORT)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	logReply := new(LogAckArg)

	err = client.Call("RaftManager.AppendLog", inArg, logReply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
}

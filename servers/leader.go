package servers

import (
	"SimpleRaft/settings"
	"log"
	"net/rpc"
	"sync"
	"errors"
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

	var cmts int

	ok_chan := make(chan bool)

	for idx := 0; idx < len(settings.AllServers); idx++ {
		ip := settings.AllServers[idx]
		if ip == this.IP {
			continue
		}
		go func(ip string) {
			var rst *LogAckArg
			rst = this.replicateLog(ip)
			for {
				if rst.Term > this.CurrentTerm { //抛出异常，由RaftManager将当前角色设置为Follower
					ok_chan <- false
					break
				}

				if !rst.Success { //因为数据不一致导致数据添加失败(更新nextIndex,然后重试)
					this.NextIndex[ip] -= 1
					rst = this.replicateLog(ip)
				} else {
					ok_chan <- true
					break
				}
			}
		}(ip)

		for {
			isOK := <- ok_chan
			if isOK {
				if cmts++;cmts>2{
					*ok = true
					return nil
				}
			}else{

			}
		}

	}



	return nil
}

func (this *Leader) replicateLog(ip string) *LogAckArg {
	nextLogIdx := this.NextIndex[ip]
	preLog := this.Logs[nextLogIdx-1]
	toSendEntries := this.Logs[nextLogIdx:]

	inArg := LogAppArg{
		Term:            this.CurrentTerm,
		LeaderID:        this.IP,
		LeaderCommitIdx: this.CommitIndex,
		PreLogIndex:     nextLogIdx - 1,
		PreLogTerm:      preLog.Term,
		Entries:         toSendEntries,
	}
	client, err := rpc.DialHTTP("tcp", ip+":"+settings.SERVERPORT)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	logReply := new(LogAckArg)

	err = client.Call("RaftManager.AppendLog", inArg, logReply)
	if err != nil {
		log.Fatalf("RaftManager AppendLog Error: NextLogIndex: %d (LogTerm: %d)",
			nextLogIdx, this.CurrentTerm)
	}

	return logReply

}

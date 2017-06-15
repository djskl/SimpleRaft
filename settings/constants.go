package settings

const TIMEOUT_MIN = 1000 	//1s
const TIMEOUT_MAX = 2000	//2s
const HEART_BEATS = 500 	//100ms
const CLIENT_WAIT = 150		//150ms
const COMMIT_WAIT = 300		//300ms
const NEWLOG_WAIT = 1000	//1000ms

const CHANN_WAIT = 1000		//1s

const RPC_WAIT = 2000		//2s

const MAJORITY = 3	//只允许5台机器

const (
	LEADER    = iota
	FOLLOWER
	CANDIDATE
)

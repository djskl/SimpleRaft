package settings

const TIMEOUT_MIN = 1000 	//300ms
const TIMEOUT_MAX = 3000 	//600ms
const HEART_BEATS = 100 	//100ms
const CLIENT_WAIT = 150		//150ms
const COMMIT_WAIT = 300		//300ms
const NEWLOG_WAIT = 1000	//1000ms

const CHANN_WAIT = 100		//100ms

const RPC_WAIT = 1000		//1s

const MAJORITY = 3	//只允许5台机器

const (
	LEADER    = iota
	FOLLOWER
	CANDIDATE
)

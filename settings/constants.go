package settings

const TIMEOUT_MIN = 1000 	//1s
const TIMEOUT_MAX = 2000	//2s

const HEART_BEATS = 300 	//300ms

const COMMIT_WAIT = 1000	//1s

const CHANN_WAIT = 1000		//1s

const RPC_WAIT = 2000		//2s

const MAJORITY = 3	//只允许5台机器

const (
	LEADER    = iota
	FOLLOWER
	CANDIDATE
)

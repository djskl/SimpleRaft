package settings

const TIMEOUT_MIN = 300 	//300ms
const TIMEOUT_MAX = 600 	//600ms
const HEART_BEATS = 100 	//100ms
const CLIENT_WAIT = 150		//150ms
const NEWLOG_WAIT = 1000	//1000ms
const COMMIT_WAIT = 300		//300ms

const CHANN_WAIT = 100		//100ms

const (
	LEADER    = iota
	FOLLOWER
	CANDIDATE
)

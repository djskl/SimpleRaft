package settings

const TIMEOUT_MIN = 300 //300ms
const TIMEOUT_MAX = 600 //600ms
const HEART_BEATS = 100 //100ms

const UUIDSIZE = 5

const (
	LEADER    = iota
	FOLLOWER
	CANDIDATE
)

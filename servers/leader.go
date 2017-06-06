package servers

type Leader struct {
	BaseRole
	NextIndex  []int
	MatchIndex []int
}

func (s *Leader) init() {
	if s.ID == ""{
		s.ID =
	}
}
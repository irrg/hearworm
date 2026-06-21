package chapter

import "time"

type Chapter struct {
	Title string
	Start time.Duration
	End   time.Duration
}

type List []Chapter

func (l List) TotalDuration() time.Duration {
	if len(l) == 0 {
		return 0
	}
	return l[len(l)-1].End
}

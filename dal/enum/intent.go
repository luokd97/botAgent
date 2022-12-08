package enum

type QueryDuration int

const (
	Yesterday   = 0
	Recent7day  = 1
	Recent30day = 2
	Recent90day = 3
	LastWeek    = 4
	LastMonth   = 5
)

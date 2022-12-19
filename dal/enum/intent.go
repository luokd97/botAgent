package enum

type QueryDuration int

const (
	Yesterday   QueryDuration = 0
	Recent7day  QueryDuration = 1
	Recent30day QueryDuration = 2
	Recent90day QueryDuration = 3
	LastWeek    QueryDuration = 4
	LastMonth   QueryDuration = 5
)

func (d QueryDuration) String() string {
	switch d {
	case Yesterday:
		return "topn_recent1day"
	case Recent7day:
		return "topn_recent7day"
	case Recent30day:
		return "topn_recent30day"
	case Recent90day:
		return "topn_recent90day"
	case LastWeek:
		return "topn_lastweek"
	case LastMonth:
		return "topn_lastmonth"
	default:
		return ""
	}
}

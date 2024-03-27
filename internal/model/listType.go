package model

type ListType int32

const (
	BlackList ListType = -1
	WhiteList ListType = 1
	NotInList ListType = 0
)

func (lt ListType) String() string {
	switch lt {
	case BlackList:
		return "BlackList"
	case WhiteList:
		return "WhiteList"
	case NotInList:
		return "NotInList"
	default:
		return "Unknown"
	}
}

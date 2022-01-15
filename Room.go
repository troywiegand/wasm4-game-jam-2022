package main

type Room struct {
	Id          int64
	Name        string
	NearbyRooms []int64
	Width       uint
	Height      uint
}

func (x Room) isEqual(y Room) bool {
	if x.Id == y.Id {
		return true
	} else {
		return false
	}
}

func (x Room) isNearby(y Room) bool {
	for _, i := range x.NearbyRooms {
		if y.Id == i {
			return true
		}
	}
	return false
}

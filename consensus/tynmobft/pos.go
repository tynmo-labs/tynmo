package tynmobft

func GetSprint(height uint64) uint64 {
	return height - height%SprintSize
}

func GetSprintRound(height uint64) uint64 {
	return height % SprintSize
}

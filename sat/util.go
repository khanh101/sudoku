package sat

func abs(x int) int {
	if x > 0 {
		return +x
	}
	if x < 0 {
		return -x
	}
	return 0
}

func sign(x int) int {
	if x > 0 {
		return +1
	}
	if x < 0 {
		return -1
	}
	return 0
}

package mock

// AddNum params add
func AddNum(a int, b int) int {
	return a + b
}

// AddAndParam params add and return
func AddAndParam(a int, b int) (int, int, int) {
	return a + b, a, b
}

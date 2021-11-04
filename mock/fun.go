package mock

func AddNum(a int, b int) int {
	return a + b
}

func AddAndParam(a int, b int) (int, int, int) {

	return a + b, a, b
}

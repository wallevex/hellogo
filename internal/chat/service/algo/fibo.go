package algo

func fiboSlow(n int64) int64 {
	if n > 50 {
		return -1
	}
	if n == 0 || n == 1 {
		return 1
	}
	return fiboSlow(n-1) + fiboSlow(n-2)
}

func Fibo(n int64) int64 {
	res := fiboSlow(n)
	return res
}

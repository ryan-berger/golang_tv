package tests

func Src(a chan int) int {
	return <-a
}

func Tgt(a chan int) int {
	return <-a
}

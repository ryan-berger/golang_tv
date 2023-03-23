package tests

func Src(a *struct{ b, c string }) int {
	if len(a.c) > len(a.b) {
		return (len(a.c) + len(a.b)) << 3
	}
	return len(a.c) + len(a.b)
}

func Tgt(a *struct{ b, c string }) int {
	return len(a.b) + len(a.c)
}

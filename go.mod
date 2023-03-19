module github.com/ryan-berger/golang_tv

go 1.20

require (
	golang.org/x/arch v0.3.0
	golang.org/x/sys v0.6.0
	tinygo.org/x/go-llvm v0.0.0-20221212185523-e80bc424a2b1
)

replace tinygo.org/x/go-llvm v0.0.0-20221212185523-e80bc424a2b1 => github.com/ryan-berger/go-llvm v0.0.0-20230319020707-41fdb91ad8a9

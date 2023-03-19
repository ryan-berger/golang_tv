// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/ryan-berger/golang_tv/alive"
	"github.com/ryan-berger/golang_tv/internal/gen"
	"os"
	"tinygo.org/x/go-llvm"
)

func main() {
	src, tgt, err := gen.GenSSA("tests/validate.go")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ctx := llvm.NewContext()
	module := ctx.NewModule("main")

	fmt.Println(src, tgt)
	srcIR, tgtIR := gen.SSA2LLVM(module, src, tgt)

	alive.Validate(module, srcIR, tgtIR)
}

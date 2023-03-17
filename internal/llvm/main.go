// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/ryan-berger/golang_tv/internal/alive"
	"tinygo.org/x/go-llvm"
)

func main() {
	llvm.InitializeAllTargets()
	llvm.InitializeAllTargetMCs()
	llvm.InitializeAllTargetInfos()
	llvm.InitializeAllAsmParsers()
	llvm.InitializeAllAsmPrinters()

	ctx := llvm.NewContext()

	module := ctx.NewModule("main")
	builder := ctx.NewBuilder()

	defer builder.Dispose()

	fnType := llvm.FunctionType(ctx.VoidType(), nil, false)

	llvmFn := llvm.AddFunction(module, "src", fnType)

	bb := ctx.AddBasicBlock(llvmFn, fmt.Sprintf("%s_bb", "src"))
	builder.SetInsertPointAtEnd(bb)

	builder.CreateAdd(llvm.ConstInt(ctx.Int64Type(), 12, false), llvm.ConstInt(ctx.Int64Type(), 33, false), "add")
	builder.CreateRetVoid()
	if err := llvm.VerifyModule(module, llvm.PrintMessageAction); err != nil {
		panic(err)
	}
	llvmFn.Dump()

	alive.Validate(module)
}

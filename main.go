// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"fmt"
	"github.com/ryan-berger/golang_tv/internal/src/cmd/compile/arm64"
	"github.com/ryan-berger/golang_tv/internal/src/cmd/compile/base"
	"github.com/ryan-berger/golang_tv/internal/src/cmd/compile/deadcode"
	"github.com/ryan-berger/golang_tv/internal/src/cmd/compile/devirtualize"
	"github.com/ryan-berger/golang_tv/internal/src/cmd/compile/escape"
	"github.com/ryan-berger/golang_tv/internal/src/cmd/compile/inline"
	"github.com/ryan-berger/golang_tv/internal/src/cmd/compile/ir"
	"github.com/ryan-berger/golang_tv/internal/src/cmd/compile/loopvar"
	"github.com/ryan-berger/golang_tv/internal/src/cmd/compile/noder"
	"github.com/ryan-berger/golang_tv/internal/src/cmd/compile/ssagen"
	"github.com/ryan-berger/golang_tv/internal/src/cmd/compile/typecheck"
	"github.com/ryan-berger/golang_tv/internal/src/cmd/compile/types"
	"github.com/ryan-berger/golang_tv/internal/src/cmd/compile/walk"
	"github.com/ryan-berger/golang_tv/internal/src/cmd/internal_compile/obj"
	"github.com/ryan-berger/golang_tv/internal/src/cmd/internal_compile/objabi"
	"os"
)

// main is a minimal example to get a single file compiled into SSA form.
// This may not be entirely complete, but is mostly copy-pasted from gc.Compile()
// - ryan-berger
func main() {
	arm64.Init(&ssagen.Arch)

	base.Ctxt = obj.Linknew(ssagen.Arch.LinkArch)
	base.Ctxt.DiagFunc = base.Errorf
	base.Ctxt.DiagFlush = base.FlushErrors
	base.Ctxt.Bso = bufio.NewWriter(os.Stdout)
	base.Ctxt.Pkgpath = "./tests"

	// set the concurrency to 1 so that Buildssa() can have a single cache
	// element to build a functon to
	// - ryan-berger
	base.Flag.LowerC = 1

	// UseBASEntries is preferred because it shaves about 2% off build time, but LLDB, dsymutil, and dwarfdump
	// on Darwin don't support it properly, especially since macOS 10.14 (Mojave).  This is exposed as a flag
	// to allow testing with LLVM tools on Linux, and to help with reporting this bug to the LLVM project.
	// See bugs 31188 and 21945 (CLs 170638, 98075, 72371).
	base.Ctxt.UseBASEntries = base.Ctxt.Headtype != objabi.Hdarwin

	symABIs := ssagen.NewSymABIs()
	if base.Flag.SymABIs != "" {
		symABIs.ReadSymABIs(base.Flag.SymABIs)
	}

	types.LocalPkg = types.NewPkg(base.Ctxt.Pkgpath, "")
	ssagen.Arch.LinkArch.Init(base.Ctxt)

	ir.IsIntrinsicCall = ssagen.IsIntrinsicCall
	inline.SSADumpInline = ssagen.DumpInline
	ssagen.InitEnv()
	ssagen.InitTables()

	types.PtrSize = ssagen.Arch.LinkArch.PtrSize
	types.RegSize = ssagen.Arch.LinkArch.RegSize
	types.MaxWidth = ssagen.Arch.MAXWIDTH

	typecheck.Target = new(ir.Package)

	typecheck.InitUniverse()
	typecheck.InitRuntime()

	// Parse and typecheck input.
	noder.LoadPackage([]string{"./tests/dummy.go"})

	ssagen.InitConfig()
	// Eliminate some obviously dead code.
	// Must happen after typechecking.
	for _, n := range typecheck.Target.Decls {
		if n.Op() == ir.ODCLFUNC {
			deadcode.Func(n.(*ir.Func))
		}
	}

	// Compute Addrtaken for names.
	// We need to wait until typechecking is done so that when we see &x[i]
	// we know that x has its address taken if x is an array, but not if x is a slice.
	// We compute Addrtaken in bulk here.
	// After this phase, we maintain Addrtaken incrementally.
	if typecheck.DirtyAddrtaken {
		typecheck.ComputeAddrtaken(typecheck.Target.Decls)
		typecheck.DirtyAddrtaken = false
	}
	typecheck.IncrementalAddrtaken = true

	noder.MakeWrappers(typecheck.Target) // must happen after inlining

	// Devirtualize and get variable capture right in for loops
	var transformed []*ir.Name
	for _, n := range typecheck.Target.Decls {
		if n.Op() == ir.ODCLFUNC {
			devirtualize.Func(n.(*ir.Func))
			transformed = append(transformed, loopvar.ForCapture(n.(*ir.Func))...)
		}
	}
	ir.CurFunc = nil

	symABIs.GenABIWrappers()
	escape.Funcs(typecheck.Target.Decls)

	var fn *ir.Func
	for _, n := range typecheck.Target.Decls {
		if n.Op() == ir.ODCLFUNC {
			fn = n.(*ir.Func)
		}
	}

	// Calculate parameter offsets.
	types.CalcSize(fn.Type())

	typecheck.DeclContext = ir.PAUTO
	ir.CurFunc = fn
	walk.Walk(fn)
	ir.CurFunc = nil // enforce no further uses of CurFunc
	typecheck.DeclContext = ir.PEXTERN

	ssad := ssagen.Buildssa(fn, 0)

	for _, b := range ssad.Blocks {
		for _, val := range b.Values {
			fmt.Println(val.Op)
		}
	}

	fmt.Println("-------")
	fmt.Println(ssad)
}

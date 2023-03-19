package gen

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
	"github.com/ryan-berger/golang_tv/internal/src/cmd/compile/liveness"
	"github.com/ryan-berger/golang_tv/internal/src/cmd/compile/loopvar"
	"github.com/ryan-berger/golang_tv/internal/src/cmd/compile/noder"
	"github.com/ryan-berger/golang_tv/internal/src/cmd/compile/reflectdata"
	"github.com/ryan-berger/golang_tv/internal/src/cmd/compile/ssa"
	"github.com/ryan-berger/golang_tv/internal/src/cmd/compile/ssagen"
	"github.com/ryan-berger/golang_tv/internal/src/cmd/compile/typecheck"
	"github.com/ryan-berger/golang_tv/internal/src/cmd/compile/types"
	"github.com/ryan-berger/golang_tv/internal/src/cmd/compile/walk"
	"github.com/ryan-berger/golang_tv/internal/src/cmd/internal_compile/obj"
	"github.com/ryan-berger/golang_tv/internal/src/cmd/internal_compile/objabi"
	"os"
	"path/filepath"
)

func ssaIt(irFn *ir.Func, worker int) *ssa.Func {
	// Calculate parameter offsets.
	types.CalcSize(irFn.Type())

	// set up some symbol lookup table, necessary for calling runtime functions
	// - ryan-berger
	ir.InitLSym(irFn, false)
	a := ssagen.AbiForBodylessFuncStackMap(irFn)
	abiInfo := a.ABIAnalyzeFuncType(irFn.Type().FuncType()) // abiInfo has spill/home locations for wrapper
	liveness.WriteFuncMap(irFn, abiInfo)

	// do the noding
	// - ryan-berger
	typecheck.DeclContext = ir.PAUTO
	ir.CurFunc = irFn
	walk.Walk(irFn)
	ir.CurFunc = nil // enforce no further uses of CurFunc
	typecheck.DeclContext = ir.PEXTERN

	return ssagen.Buildssa(irFn, worker)
}

// GenSSA is a minimal example to get a single file compiled into SSA form.
// This may not be entirely complete, but is mostly copy-pasted from gc.Compile()
// TODO: it would be nice if we could ask nicely to return the functions we want rather than hard-coding Source and Target
//
// TODO: allow for easier pass -> pass validation on a single function doing fn -> SSA[Unoptimized], SSA[Unoptimized] -> SSA[Optimized]
//
//	this will likely require more open "pass" editing
//
// - ryan-berger
func GenSSA(filePath string) (*ssa.Func, *ssa.Func, error) {
	pkgPath, _ := filepath.Split(filePath)

	arm64.Init(&ssagen.Arch)
	base.Ctxt = obj.Linknew(ssagen.Arch.LinkArch)
	base.Ctxt.DiagFunc = base.Errorf
	base.Ctxt.DiagFlush = base.FlushErrors
	base.Ctxt.Bso = bufio.NewWriter(os.Stdout)

	base.Ctxt.Pkgpath = pkgPath

	// set the concurrency to 2 so that Buildssa() can compile two functions withougt
	// overwriting eachother
	// TODO: understand Buildssa better
	//- ryan-berger
	base.Flag.LowerC = 2

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

	// pseudo-package, for scoping
	types.BuiltinPkg = types.NewPkg("go.builtin", "") // TODO(gri) name this package go.builtin?
	types.BuiltinPkg.Prefix = "go:builtin"

	// pseudo-package, accessed by import "unsafe"
	types.UnsafePkg = types.NewPkg("unsafe", "unsafe")

	// Pseudo-package that contains the compiler's builtin
	// declarations for package runtime. These are declared in a
	// separate package to avoid conflicts with package runtime's
	// actual declarations, which may differ intentionally but
	// insignificantly.
	ir.Pkgs.Runtime = types.NewPkg("go.runtime", "runtime")
	ir.Pkgs.Runtime.Prefix = "runtime"

	// pseudo-packages used in symbol tables
	ir.Pkgs.Itab = types.NewPkg("go.itab", "go.itab")
	ir.Pkgs.Itab.Prefix = "go:itab"

	// pseudo-package used for methods with anonymous receivers
	ir.Pkgs.Go = types.NewPkg("go", "")

	ssagen.Arch.LinkArch.Init(base.Ctxt)

	ir.IsIntrinsicCall = ssagen.IsIntrinsicCall
	inline.SSADumpInline = ssagen.DumpInline
	ssagen.InitEnv()
	ssagen.InitTables()

	types.PtrSize = ssagen.Arch.LinkArch.PtrSize
	types.RegSize = ssagen.Arch.LinkArch.RegSize
	types.MaxWidth = ssagen.Arch.MAXWIDTH

	typecheck.Target = new(ir.Package)

	typecheck.NeedRuntimeType = reflectdata.NeedRuntimeType // TODO(rsc): TypeSym for lock?

	typecheck.InitUniverse()
	typecheck.InitRuntime()

	// Parse and typecheck input.
	noder.LoadPackage([]string{filePath})

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

	var srcIR, tgtIR *ir.Func
	for _, n := range typecheck.Target.Decls {
		if n.Op() == ir.ODCLFUNC {
			switch f := n.(*ir.Func); ir.FuncName(f) {
			case "Src":
				srcIR = f
			case "Tgt":
				tgtIR = f
			}
		}
	}

	if srcIR == nil || tgtIR == nil {
		return nil, nil, fmt.Errorf("Must have a func named Src and Tgt")
	}

	return ssaIt(srcIR, 0), ssaIt(tgtIR, 1), nil
}

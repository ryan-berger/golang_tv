package gen

import (
	"github.com/ryan-berger/golang_tv/alive"
	"github.com/ryan-berger/golang_tv/internal/src/cmd/compile/ssa"
	"github.com/ryan-berger/golang_tv/internal/src/cmd/compile/types"
	"math"
	"tinygo.org/x/go-llvm"
)

func init() {
	llvm.InitializeAllTargets()
	llvm.InitializeAllTargetMCs()
	llvm.InitializeAllTargetInfos()
	llvm.InitializeAllAsmParsers()
	llvm.InitializeAllAsmPrinters()
}

func (l *llvmGenerator) goTypeToLLVMType(typ *types.Type) llvm.Type {
	switch typ.Kind() {
	case types.TSTRUCT:
		var types []llvm.Type
		for _, f := range typ.Fields().Slice() {
			types = append(types, l.goTypeToLLVMType(f.Type))
		}
		return l.ctx.StructType(types, false)
	case types.TINT: // TODO: actually resolve int type correctly
		return l.ctx.Int64Type()
	case types.TFLOAT64:
		return l.ctx.FloatType()
	case types.TFLOAT32:
		return l.ctx.DoubleType()
	case types.TFUNC:
		var argTypes []llvm.Type

		fn := typ.FuncType()
		fields := fn.Params.Fields()

		for _, field := range fields.Slice() {
			argTypes = append(argTypes, l.goTypeToLLVMType(field.Type))
		}

		return llvm.FunctionType(
			l.goTypeToLLVMType(fn.Results),
			[]llvm.Type{l.goTypeToLLVMType(fn.Params)}, false)
	}

	return l.ctx.VoidType()
}

func (l *llvmGenerator) genVal(fn llvm.Value, typ llvm.Type, val *ssa.Value) {
	var v llvm.Value
	switch val.Op {
	case ssa.OpConst64F:
		v = llvm.ConstFloat(l.ctx.FloatType(), math.Float64frombits(uint64(val.AuxInt)))
	case ssa.OpConst64:
		v = llvm.ConstInt(l.ctx.Int64Type(), uint64(val.AuxInt), true)
	case ssa.OpLeq64:
		v = l.builder.CreateICmp(llvm.IntSLE, l.curFn[val.Args[0].String()], l.curFn[val.Args[1].String()], val.String())
	case ssa.OpArgIntReg, ssa.OpArgFloatReg:
		v = l.builder.CreateExtractValue(fn.Param(0), int(val.AuxInt), val.String())
	case ssa.OpAdd64:
		v = l.builder.CreateAdd(l.curFn[val.Args[0].String()], l.curFn[val.Args[1].String()], val.String())
	case ssa.OpAdd64F:
		v = l.builder.CreateFAdd(l.curFn[val.Args[0].String()], l.curFn[val.Args[1].String()], val.String())
	case ssa.OpSub64:
		v = l.builder.CreateSub(l.curFn[val.Args[0].String()], l.curFn[val.Args[1].String()], val.String())
	case ssa.OpRsh64x64:
		rhs := l.curFn[val.Args[1].String()]

		shouldClamp := l.builder.CreateICmp(llvm.IntSGE,
			rhs,
			llvm.ConstInt(l.ctx.Int64Type(), 63, false), "should_clamp")
		shift := l.builder.CreateSelect(shouldClamp,
			llvm.ConstInt(l.ctx.Int64Type(), 63, false), rhs, "clamped")

		v = l.builder.CreateAShr(l.curFn[val.Args[0].String()], shift, val.String())

	case ssa.OpLsh64x64:
		rhs := l.curFn[val.Args[1].String()]
		// clamp
		shouldClamp := l.builder.CreateICmp(llvm.IntSGE,
			rhs,
			llvm.ConstInt(l.ctx.Int64Type(), 63, false), "should_clamp")

		shift := l.builder.CreateSelect(shouldClamp,
			llvm.ConstInt(l.ctx.Int64Type(), 63, false), rhs, "clamped")

		v = l.builder.CreateLShr(l.curFn[val.Args[0].String()], shift, val.String())
	case ssa.OpMakeResult:
		iv := l.builder.CreateInsertValue(llvm.Undef(typ.ReturnType()), l.curFn[val.Args[0].String()], 0, "first")
		l.builder.CreateRet(iv)
		return
	}

	l.curFn[val.String()] = v
}

func (l *llvmGenerator) genFn(fn *ssa.Func) llvm.Value {
	typ := l.goTypeToLLVMType(fn.Type)

	name := fn.Name
	switch name {
	case "Src":
		name = "src"
	case "Tgt":
		name = "tgt"
	}

	llvmFn := llvm.AddFunction(l.module, name, typ)

	indexedBlocks := make(map[string]llvm.BasicBlock, len(fn.Blocks))
	for _, b := range fn.Blocks {
		if b.Kind == ssa.BlockInvalid { // don't index invalid blocks!
			continue
		}
		indexedBlocks[b.String()] = l.ctx.AddBasicBlock(llvmFn, b.String())
	}

	// TODO: walk this in correct order
	for _, b := range fn.Blocks {
		if b.Kind == ssa.BlockInvalid { // skip over invalid block code generation!
			continue
		}

		bb := indexedBlocks[b.String()]
		l.builder.SetInsertPointAtEnd(bb)

		for _, v := range b.Values {
			l.genVal(llvmFn, typ, v)
		}

		switch b.Kind {
		case ssa.BlockIf:
			cond := l.curFn[b.Controls[0].String()]
			then, otherwise := indexedBlocks[b.Succs[0].Block().String()], indexedBlocks[b.Succs[1].Block().String()]
			l.builder.CreateCondBr(cond, then, otherwise)
		case ssa.BlockExit:
			l.builder.CreateUnreachable()
		}
	}

	return llvmFn
}

type llvmGenerator struct {
	ctx     llvm.Context
	builder llvm.Builder
	module  llvm.Module
	curFn   map[string]llvm.Value
}

func SSA2LLVM(module llvm.Module, src, tgt *ssa.Func) (llvm.Value, llvm.Value) {
	ctx := module.Context()
	builder := ctx.NewBuilder()

	defer builder.Dispose()

	b := &llvmGenerator{
		ctx:     ctx,
		builder: builder,
		module:  module,
		curFn:   make(map[string]llvm.Value),
	}

	srcLLVM := b.genFn(src)
	//_ = b.genFn(src)
	b.curFn = make(map[string]llvm.Value)
	tgtLLVM := b.genFn(tgt)
	//_ = b.genFn(tgt)

	if err := llvm.VerifyModule(module, llvm.PrintMessageAction); err != nil {
		panic(err)
	}

	alive.Validate(module, srcLLVM, tgtLLVM)
	return srcLLVM, tgtLLVM
}

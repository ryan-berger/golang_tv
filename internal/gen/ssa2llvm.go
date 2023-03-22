package gen

import (
	"fmt"
	"github.com/ryan-berger/golang_tv/internal/src/cmd/compile/ir"
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

		return llvm.FunctionType(l.goTypeToLLVMType(fn.Results), argTypes, false)
	case types.TSLICE:
		return l.ctx.StructType([]llvm.Type{
			llvm.PointerType(l.goTypeToLLVMType(typ.Elem()), 0),
			l.ctx.Int64Type(),
			l.ctx.Int64Type(),
		}, false)
	case types.TARRAY:
		return llvm.ArrayType(
			l.goTypeToLLVMType(typ.Elem()),
			int(typ.NumElem()),
		)
	case types.TPTR:
		return llvm.PointerType(l.goTypeToLLVMType(typ.Underlying()), 0)
	default:
		panic(fmt.Sprintf(`typecheck err: verification with type "%s" is not currently supported by golang_tv`, typ.String()))

	}

	return l.ctx.VoidType()
}

var icmpMap = map[ssa.Op]llvm.IntPredicate{
	// lt
	ssa.OpLess64: llvm.IntSLT,
	ssa.OpLess32: llvm.IntSLT,
	ssa.OpLess16: llvm.IntSLT,
	ssa.OpLess8:  llvm.IntSLT,

	// le
	ssa.OpLeq64: llvm.IntSLE,
	ssa.OpLeq32: llvm.IntSLE,
	ssa.OpLeq16: llvm.IntSLE,
	ssa.OpLeq8:  llvm.IntSLE,

	// eq
	ssa.OpEq64: llvm.IntEQ,
	ssa.OpEq32: llvm.IntEQ,
	ssa.OpEq16: llvm.IntEQ,
	ssa.OpEq8:  llvm.IntEQ,

	// ne
	ssa.OpNeq64: llvm.IntNE,
	ssa.OpNeq32: llvm.IntNE,
	ssa.OpNeq16: llvm.IntNE,
	ssa.OpNeq8:  llvm.IntNE,

	// ult
	ssa.OpLess64U: llvm.IntSLT,
	ssa.OpLess32U: llvm.IntSLT,
	ssa.OpLess16U: llvm.IntSLT,
	ssa.OpLess8U:  llvm.IntSLT,

	// ule
	ssa.OpLeq64U: llvm.IntSLE,
	ssa.OpLeq32U: llvm.IntSLE,
	ssa.OpLeq16U: llvm.IntSLE,
	ssa.OpLeq8U:  llvm.IntSLE,
}

func (l *llvmGenerator) genVal(fn llvm.Value, val *ssa.Value) {
	var v llvm.Value
	switch val.Op {
	case ssa.OpConst64F:
		v = llvm.ConstFloat(l.ctx.FloatType(), math.Float64frombits(uint64(val.AuxInt)))
	case ssa.OpConst64:
		v = llvm.ConstInt(l.ctx.Int64Type(), uint64(val.AuxInt), true)
	case ssa.OpLoad:
		llvmType := l.goTypeToLLVMType(val.Args[0].Type.Elem())

		arg := l.curFn[val.Args[0].String()]

		v = l.builder.CreateLoad(llvmType, arg, val.String())
	case ssa.OpLocalAddr:
		offset := val.Aux.(*ir.Name)
		n := l.curFn[offset.Sym().Name]

		a := l.builder.CreateAlloca(l.ctx.Int8Type(), fmt.Sprintf("%s_ptr", offset.Sym().Name))
		l.builder.CreateStore(n, a)

		v = a
	case ssa.OpOffPtr:
		arg := l.curFn[val.Args[0].String()]

		// calculate GEP as if the ptr was a char*, opaque pointers make this easy!
		charPtr := llvm.PointerType(l.ctx.Int8Type(), 0)
		v = l.builder.CreateGEP(charPtr, arg,
			[]llvm.Value{
				llvm.ConstInt(l.ctx.Int64Type(), uint64(val.AuxInt), false),
			},
			val.String())
	case ssa.OpIsInBounds:
		index := l.curFn[val.Args[0].String()]
		length := l.curFn[val.Args[1].String()]

		// make sure that the int is greather t
		gtZero := l.builder.CreateICmp(llvm.IntSGE, index, llvm.ConstInt(l.ctx.Int64Type(), 0, false), "gtZero")
		ltLength := l.builder.CreateICmp(llvm.IntSLT, index, length, "ltLength")

		v = l.builder.CreateAnd(gtZero, ltLength, val.String())
	case ssa.OpCom64, ssa.OpCom32, ssa.OpCom16, ssa.OpCom8:
		v = l.builder.CreateNot(l.curFn[val.Args[0].String()], val.String())
	case ssa.OpAnd64:
		v = l.builder.CreateAnd(l.curFn[val.Args[0].String()], l.curFn[val.Args[1].String()], val.String())
	case ssa.OpLess64, ssa.OpLess32, ssa.OpLess16, ssa.OpLess8,
		ssa.OpLeq64, ssa.OpLeq32, ssa.OpLeq16, ssa.OpLeq8,
		ssa.OpEq64, ssa.OpEq32, ssa.OpEq16, ssa.OpEq8,
		ssa.OpNeq64, ssa.OpNeq32, ssa.OpNeq16, ssa.OpNeq8,
		ssa.OpLess64U, ssa.OpLess32U, ssa.OpLess16U, ssa.OpLess8U,
		ssa.OpLeq64U, ssa.OpLeq32U, ssa.OpLeq16U, ssa.OpLeq8U:
		v = l.builder.CreateICmp(icmpMap[val.Op], l.curFn[val.Args[0].String()], l.curFn[val.Args[1].String()], val.String())
	case ssa.OpArgIntReg, ssa.OpArgFloatReg:
		switch aux := val.Aux.(type) {
		case *ssa.AuxNameOffset:
			param := l.curFn[aux.Name.Sym().Name]
			if param.Type().TypeKind() == llvm.StructTypeKind {
				v = l.builder.CreateExtractValue(l.curFn[aux.Name.Sym().Name], int(val.AuxInt), val.String())
			} else {
				v = param
			}
		default:
			panic("other aux types not yet supported")
		}
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

		v = l.builder.CreateLShr(l.curFn[val.Args[0].String()], shift, val.String())
	case ssa.OpLsh64x64:
		rhs := l.curFn[val.Args[1].String()]
		// clamp
		shouldClamp := l.builder.CreateICmp(llvm.IntSGE,
			rhs,
			llvm.ConstInt(l.ctx.Int64Type(), 63, false), "should_clamp")
		shift := l.builder.CreateSelect(shouldClamp,
			llvm.ConstInt(l.ctx.Int64Type(), 63, false), rhs, "clamped")

		v = l.builder.CreateShl(l.curFn[val.Args[0].String()], shift, val.String())
	case ssa.OpMakeResult:
		ret := fn.FunctionType().ReturnType()
		iv := l.builder.CreateInsertValue(llvm.Undef(ret), l.curFn[val.Args[0].String()], 0, "first")
		l.builder.CreateRet(iv)
		return
	default:
		panic(fmt.Sprintf("SSA op: %s is not currently supported for verification by golang_tv", val.Op.String()))
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

	goFnType := fn.Type.FuncType()
	goParams := goFnType.Params.Fields()
	for i, p := range goParams.Slice() {
		l.curFn[p.Sym.Name] = llvmFn.Param(i)
	}

	indexedBlocks := make(map[string]llvm.BasicBlock, len(fn.Blocks))
	for _, b := range fn.Blocks {
		indexedBlocks[b.String()] = l.ctx.AddBasicBlock(llvmFn, b.String())
	}

	for _, b := range fn.Blocks {
		bb := indexedBlocks[b.String()]
		l.builder.SetInsertPointAtEnd(bb)

		for _, v := range b.Values {
			l.genVal(llvmFn, v)
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

	module.Dump()
	if err := llvm.VerifyModule(module, llvm.PrintMessageAction); err != nil {
		panic(err)
	}

	return srcLLVM, tgtLLVM
}

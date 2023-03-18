package alive

// #cgo CXXFLAGS: -g -I/home/ryan/github.com/AliveToolkit/alive2 -I/home/ryan/github.com/llvm/llvm-project/llvm/include -I/home/ryan/github.com/llvm/llvm-project/build/include -std=c++20 -D_GNU_SOURCE -D_DEBUG -D_GLIBCXX_ASSERTIONS -D_LIBCPP_ENABLE_ASSERTIONS -D__STDC_CONSTANT_MACROS -D__STDC_FORMAT_MACROS -D__STDC_LIMIT_MACROS
// #cgo LDFLAGS: -L/home/ryan/github.com/AliveToolkit/alive2/build -lllvm_util -lz3 -lir -lsmt -ltools -lutil -L/home/ryan/github.com/llvm/llvm-project/build/lib -lLLVMWindowsManifest -lLLVMXRay -lLLVMLibDriver -lLLVMDlltoolDriver -lLLVMCoverage -lLLVMLineEditor -lLLVMX86TargetMCA -lLLVMX86Disassembler -lLLVMX86AsmParser -lLLVMX86CodeGen -lLLVMX86Desc -lLLVMX86Info -lLLVMOrcJIT -lLLVMWindowsDriver -lLLVMMCJIT -lLLVMJITLink -lLLVMInterpreter -lLLVMExecutionEngine -lLLVMRuntimeDyld -lLLVMOrcTargetProcess -lLLVMOrcShared -lLLVMDWP -lLLVMDebugInfoLogicalView -lLLVMDebugInfoGSYM -lLLVMOption -lLLVMObjectYAML -lLLVMObjCopy -lLLVMMCA -lLLVMMCDisassembler -lLLVMLTO -lLLVMPasses -lLLVMCFGuard -lLLVMCoroutines -lLLVMipo -lLLVMVectorize -lLLVMLinker -lLLVMInstrumentation -lLLVMFrontendOpenMP -lLLVMFrontendOpenACC -lLLVMFrontendHLSL -lLLVMExtensions -lLLVMDWARFLinkerParallel -lLLVMDWARFLinker -lLLVMGlobalISel -lLLVMMIRParser -lLLVMAsmPrinter -lLLVMSelectionDAG -lLLVMCodeGen -lLLVMObjCARCOpts -lLLVMIRPrinter -lLLVMInterfaceStub -lLLVMFileCheck -lLLVMFuzzMutate -lLLVMTarget -lLLVMScalarOpts -lLLVMInstCombine -lLLVMAggressiveInstCombine -lLLVMTransformUtils -lLLVMBitWriter -lLLVMAnalysis -lLLVMProfileData -lLLVMSymbolize -lLLVMDebugInfoPDB -lLLVMDebugInfoMSF -lLLVMDebugInfoDWARF -lLLVMObject -lLLVMTextAPI -lLLVMMCParser -lLLVMIRReader -lLLVMAsmParser -lLLVMMC -lLLVMDebugInfoCodeView -lLLVMBitReader -lLLVMFuzzerCLI -lLLVMCore -lLLVMRemarks -lLLVMBitstreamReader -lLLVMBinaryFormat -lLLVMTargetParser -lLLVMTableGen -lLLVMSupport -lLLVMDemangle -Wl,-rpath=/home/ryan/github.com/llvm/llvm-project/build/lib:/home/ryan/github.com/AliveToolkit/alive2/build
// #include "alive.h"
import "C"
import (
	"fmt"
	"tinygo.org/x/go-llvm"
	"unsafe"
)

func Validate(module llvm.Module, src, tgt llvm.Value) {
	if src.IsAFunction().IsNil() || tgt.IsAFunction().IsNil() {
		panic("src, tgt must be functions")
	}
	C.validate((C.LLVMModuleRef)(unsafe.Pointer(module.C)),
		(C.LLVMValueRef)(unsafe.Pointer(src.C)),
		(C.LLVMValueRef)(unsafe.Pointer(tgt.C)))
	fmt.Println("done!")
}
